package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/pixel-game/internal/auth"
	"github.com/yourusername/pixel-game/internal/domain"
	"github.com/yourusername/pixel-game/internal/game/effects"
	"github.com/yourusername/pixel-game/internal/game/ai"
	"github.com/yourusername/pixel-game/internal/middleware"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameRepo       domain.GameRepository
	cardRepo       domain.CardRepository
	userRepo       domain.UserRepository
	jwtManager     *auth.JWTManager
	effectExecutor *effects.Executor
	aiManager      *ai.AIManager
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameRepo domain.GameRepository, cardRepo domain.CardRepository, userRepo domain.UserRepository, jwtManager *auth.JWTManager) *GameHandler {
	return &GameHandler{
		gameRepo:       gameRepo,
		cardRepo:       cardRepo,
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		effectExecutor: effects.NewExecutor(),
		aiManager:      ai.NewAIManager(),
	}
}

// RegisterRoutes registers game routes
func (h *GameHandler) RegisterRoutes(router *gin.RouterGroup) {
	games := router.Group("/games")
	games.Use(middleware.AuthMiddleware(h.jwtManager))
	{
		games.POST("/start", h.StartGame)
		games.GET("/current", h.GetCurrentGame)
		games.GET("/:id", h.GetGame)
		games.POST("/:id/actions", h.PlayAction)
		games.POST("/:id/end-turn", h.EndTurn)
		games.POST("/:id/surrender", h.SurrenderGame)
		games.GET("/stats", h.GetGameStats)
	}
}

// StartGameRequest represents a request to start a new game
type StartGameRequest struct {
	GameMode domain.GameMode `json:"game_mode" binding:"required"`
	DeckID   *int            `json:"deck_id"`
}

// StartGame godoc
// @Summary 게임 시작
// @Description 새로운 게임 세션을 시작합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body StartGameRequest true "게임 시작 요청"
// @Success 201 {object} map[string]interface{} "생성된 게임 세션"
// @Failure 400 {object} map[string]interface{} "잘못된 요청"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 409 {object} map[string]interface{} "이미 진행 중인 게임이 있음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/start [post]
func (h *GameHandler) StartGame(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	var req StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청입니다",
		})
		return
	}

	// Check if user already has an active game
	activeGame, err := h.gameRepo.GetActiveSession(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 확인할 수 없습니다",
		})
		return
	}

	if activeGame != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "이미 진행 중인 게임이 있습니다",
			"game_id": activeGame.ID,
		})
		return
	}

	// Get user's active deck
	var deck *domain.Deck
	if req.DeckID != nil {
		deck, err = h.cardRepo.GetDeck(*req.DeckID)
		if err != nil || deck == nil || deck.UserID != userID.(int) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "유효하지 않은 덱입니다",
			})
			return
		}
	} else {
		deck, err = h.cardRepo.GetActiveDeck(userID.(int))
		if err != nil || deck == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "활성화된 덱이 없습니다",
			})
			return
		}
	}

	// Initialize game session
	session := &domain.GameSession{
		UserID:        userID.(int),
		Status:        domain.GameStatusActive,
		GameMode:      req.GameMode,
		CurrentFloor:  1,
		CurrentTurn:   1,
		TurnPhase:     domain.TurnPhaseStart,
		DeckSnapshot:  deck.CardIDs,
		TurnTimeLimit: 120, // 2 minutes per turn
	}

	// Initialize player state
	playerState := &domain.PlayerState{
		Health:       100,
		MaxHealth:    100,
		Shield:       0,
		Energy:       3,
		MaxEnergy:    3,
		Hand:         []string{},
		DrawPile:     make([]string, len(deck.CardIDs)),
		DiscardPile:  []string{},
		ExhaustPile:  []string{},
		ActivePowers: make(map[string]domain.PowerState),
		Buffs:        []domain.BuffState{},
		Debuffs:      []domain.DebuffState{},
	}
	copy(playerState.DrawPile, deck.CardIDs)
	// TODO: Shuffle draw pile

	// Draw initial hand
	playerState.DrawCards(5)

	// Initialize enemy for first floor
	enemyState := h.generateEnemy(1, req.GameMode)

	// Initialize game state
	gameState := &domain.GameState{
		FloorType:   "COMBAT",
		FloorData:   map[string]interface{}{},
		Relics:      []string{},
		Gold:        50,
		PotionSlots: 3,
		Potions:     []string{},
		CardRewards: []string{},
		Path:        h.generatePath(req.GameMode),
		CurrentNodeID: "1-1",
	}

	// Marshal states to JSON
	playerJSON, _ := json.Marshal(playerState)
	enemyJSON, _ := json.Marshal(enemyState)
	gameJSON, _ := json.Marshal(gameState)
	
	session.PlayerState = playerJSON
	session.EnemyState = enemyJSON
	session.GameState = gameJSON

	// Create session
	if err := h.gameRepo.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임을 시작할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session_id": session.ID,
		"status": session.Status,
		"game_mode": session.GameMode,
		"current_floor": session.CurrentFloor,
		"current_turn": session.CurrentTurn,
		"turn_phase": session.TurnPhase,
		"player_state": playerState,
		"enemy_state": enemyState,
		"game_state": gameState,
	})
}

// GetCurrentGame godoc
// @Summary 현재 게임 조회
// @Description 현재 진행 중인 게임 세션을 조회합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "게임 세션 정보"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 404 {object} map[string]interface{} "진행 중인 게임 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/current [get]
func (h *GameHandler) GetCurrentGame(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	session, err := h.gameRepo.GetActiveSession(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 조회할 수 없습니다",
		})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "진행 중인 게임이 없습니다",
		})
		return
	}

	playerState, enemyState, gameState, err := h.gameRepo.LoadGameState(session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 불러올 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"status": session.Status,
		"game_mode": session.GameMode,
		"current_floor": session.CurrentFloor,
		"current_turn": session.CurrentTurn,
		"turn_phase": session.TurnPhase,
		"score": session.Score,
		"player_state": playerState,
		"enemy_state": enemyState,
		"game_state": gameState,
		"last_action_at": session.LastActionAt,
	})
}

// GetGame godoc
// @Summary 게임 세션 조회
// @Description 특정 게임 세션의 정보를 조회합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "게임 세션 ID"
// @Success 200 {object} map[string]interface{} "게임 세션 정보"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "게임을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/{id} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 게임 ID입니다",
		})
		return
	}

	session, err := h.gameRepo.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임을 조회할 수 없습니다",
		})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "게임을 찾을 수 없습니다",
		})
		return
	}

	if session.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 게임에 접근할 권한이 없습니다",
		})
		return
	}

	playerState, enemyState, gameState, err := h.gameRepo.LoadGameState(session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 불러올 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"status": session.Status,
		"game_mode": session.GameMode,
		"current_floor": session.CurrentFloor,
		"current_turn": session.CurrentTurn,
		"turn_phase": session.TurnPhase,
		"score": session.Score,
		"cards_played": session.CardsPlayed,
		"damage_dealt": session.DamageDealt,
		"damage_taken": session.DamageTaken,
		"player_state": playerState,
		"enemy_state": enemyState,
		"game_state": gameState,
		"started_at": session.StartedAt,
		"completed_at": session.CompletedAt,
		"last_action_at": session.LastActionAt,
	})
}

// PlayActionRequest represents a game action request
type PlayActionRequest struct {
	ActionType domain.ActionType `json:"action_type" binding:"required"`
	CardID     *string           `json:"card_id,omitempty"`
	TargetID   *string           `json:"target_id,omitempty"`
	ActionData json.RawMessage   `json:"action_data,omitempty"`
}

// PlayAction godoc
// @Summary 게임 액션 실행
// @Description 카드 플레이, 포션 사용 등의 게임 액션을 실행합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "게임 세션 ID"
// @Param request body PlayActionRequest true "게임 액션 요청"
// @Success 200 {object} map[string]interface{} "액션 결과"
// @Failure 400 {object} map[string]interface{} "잘못된 요청"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "게임을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/{id}/actions [post]
func (h *GameHandler) PlayAction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 게임 ID입니다",
		})
		return
	}

	var req PlayActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청입니다",
		})
		return
	}

	// Get session
	session, err := h.gameRepo.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "게임을 찾을 수 없습니다",
		})
		return
	}

	if session.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 게임에 접근할 권한이 없습니다",
		})
		return
	}

	if !session.CanTakeAction() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "현재 액션을 수행할 수 없는 상태입니다",
		})
		return
	}

	// Load game state
	playerState, enemyState, gameState, err := h.gameRepo.LoadGameState(session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 불러올 수 없습니다",
		})
		return
	}

	// Process action based on type
	var result map[string]interface{}
	switch req.ActionType {
	case domain.ActionTypePlayCard:
		result, err = h.processPlayCard(session, playerState, enemyState, gameState, req.CardID, req.TargetID)
	case domain.ActionTypeUsePotion:
		result, err = h.processUsePotion(session, playerState, enemyState, gameState, req.ActionData)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "지원하지 않는 액션 타입입니다",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Record action
	action := &domain.GameAction{
		SessionID:  sessionID,
		ActionType: string(req.ActionType),
		CardID:     req.CardID,
		TargetID:   req.TargetID,
		ActionData: req.ActionData,
	}
	h.gameRepo.RecordAction(action)

	// Update session statistics
	if req.ActionType == domain.ActionTypePlayCard {
		session.CardsPlayed++
	}

	// Save updated game state
	if err := h.gameRepo.SaveGameState(sessionID, playerState, enemyState, gameState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 저장할 수 없습니다",
		})
		return
	}

	// Update session
	if err := h.gameRepo.UpdateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 세션을 업데이트할 수 없습니다",
		})
		return
	}

	result["player_state"] = playerState
	result["enemy_state"] = enemyState
	result["game_state"] = gameState

	c.JSON(http.StatusOK, result)
}

// EndTurn godoc
// @Summary 턴 종료
// @Description 현재 턴을 종료하고 다음 턴으로 진행합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "게임 세션 ID"
// @Success 200 {object} map[string]interface{} "턴 종료 결과"
// @Failure 400 {object} map[string]interface{} "잘못된 요청"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "게임을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/{id}/end-turn [post]
func (h *GameHandler) EndTurn(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 게임 ID입니다",
		})
		return
	}

	// Get session
	session, err := h.gameRepo.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "게임을 찾을 수 없습니다",
		})
		return
	}

	if session.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 게임에 접근할 권한이 없습니다",
		})
		return
	}

	if session.TurnPhase != domain.TurnPhaseMain {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "현재 턴을 종료할 수 없는 상태입니다",
		})
		return
	}

	// Load game state
	playerState, enemyState, gameState, err := h.gameRepo.LoadGameState(session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 불러올 수 없습니다",
		})
		return
	}

	// Process end turn
	// 1. Move hand cards to discard pile
	playerState.DiscardPile = append(playerState.DiscardPile, playerState.Hand...)
	playerState.Hand = []string{}

	// 2. Enemy turn
	session.TurnPhase = domain.TurnPhaseEnemy
	enemyActions := h.processEnemyTurn(session, playerState, enemyState, gameState)

	// Check if player is defeated
	if playerState.Health <= 0 {
		session.Status = domain.GameStatusFailed
		h.gameRepo.EndSession(sessionID, domain.GameStatusFailed)
		c.JSON(http.StatusOK, gin.H{
			"message": "게임 오버",
			"result": "defeat",
			"enemy_actions": enemyActions,
			"player_state": playerState,
			"enemy_state": enemyState,
			"game_state": gameState,
		})
		return
	}

	// Check if enemy is defeated
	if enemyState.Health <= 0 {
		// Process victory
		result := h.processVictory(session, playerState, enemyState, gameState)
		c.JSON(http.StatusOK, result)
		return
	}

	// 3. Start new turn
	session.CurrentTurn++
	session.TurnPhase = domain.TurnPhaseStart
	
	// Reset energy
	playerState.Energy = playerState.MaxEnergy
	
	// Draw cards for new turn
	playerState.DrawCards(5)

	// Update buffs/debuffs duration
	h.updateEffectDurations(playerState, enemyState)

	// Save state
	if err := h.gameRepo.SaveGameState(sessionID, playerState, enemyState, gameState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 상태를 저장할 수 없습니다",
		})
		return
	}

	session.TurnPhase = domain.TurnPhaseMain
	if err := h.gameRepo.UpdateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임 세션을 업데이트할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "턴 종료",
		"current_turn": session.CurrentTurn,
		"enemy_actions": enemyActions,
		"player_state": playerState,
		"enemy_state": enemyState,
		"game_state": gameState,
	})
}

// SurrenderGame godoc
// @Summary 게임 포기
// @Description 현재 게임을 포기합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "게임 세션 ID"
// @Success 200 {object} map[string]interface{} "게임 포기 결과"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "게임을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/{id}/surrender [post]
func (h *GameHandler) SurrenderGame(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 게임 ID입니다",
		})
		return
	}

	session, err := h.gameRepo.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "게임을 찾을 수 없습니다",
		})
		return
	}

	if session.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 게임에 접근할 권한이 없습니다",
		})
		return
	}

	if session.Status != domain.GameStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "이미 종료된 게임입니다",
		})
		return
	}

	// End session
	if err := h.gameRepo.EndSession(sessionID, domain.GameStatusFailed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "게임을 종료할 수 없습니다",
		})
		return
	}

	// Update user stats
	if err := h.userRepo.IncrementGamesPlayed(userID.(int)); err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "게임을 포기했습니다",
		"session_id": sessionID,
		"final_score": session.Score,
		"final_floor": session.CurrentFloor,
	})
}

// GetGameStats godoc
// @Summary 게임 통계 조회
// @Description 사용자의 게임 통계를 조회합니다
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.UserGameStats "게임 통계"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/games/stats [get]
func (h *GameHandler) GetGameStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	stats, err := h.gameRepo.GetUserGameStats(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "통계를 조회할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Helper methods

func (h *GameHandler) generateEnemy(floor int, gameMode domain.GameMode) *domain.EnemyState {
	// 층수에 따른 적 타입 결정
	var enemyType, enemyName, aiType string
	baseHealth := 40
	
	switch {
	case floor <= 2:
		enemyType = "BASIC_ENEMY"
		enemyName = "사이버 드론"
		aiType = "balanced"
		baseHealth = 40
	case floor <= 4:
		enemyType = "BRUTE"
		enemyName = "사이버 워리어"
		aiType = "aggressive"
		baseHealth = 60
	case floor <= 6:
		enemyType = "GUARDIAN"
		enemyName = "사이버 가디언"
		aiType = "defensive"
		baseHealth = 80
	case floor%3 == 0: // 보스급
		enemyType = "ELITE"
		enemyName = "사이버 로드"
		aiType = "balanced"
		baseHealth = 120
	default:
		enemyType = "BASIC_ENEMY"
		enemyName = "사이버 스컬지"
		aiType = "balanced"
		baseHealth = 50
	}
	
	// 체력 계산 (층수에 따라 증가)
	maxHealth := baseHealth + (floor * 8)
	
	// 적 상태 생성
	enemy := &domain.EnemyState{
		ID:           fmt.Sprintf("enemy_%d_%s", floor, enemyType),
		Name:         enemyName,
		Health:       maxHealth,
		MaxHealth:    maxHealth,
		Shield:       0,
		ActivePowers: []domain.PowerState{},
		Buffs:        []domain.BuffState{},
		Debuffs:      []domain.DebuffState{},
	}
	
	// AI 시스템을 사용해서 첫 번째 의도 계산
	intent, err := h.generateInitialIntent(enemy, aiType, floor)
	if err != nil {
		// 에러 시 기본 의도 설정
		enemy.Intent = domain.EnemyIntent{
			Type:        "ATTACK",
			Value:       10 + floor,
			Description: fmt.Sprintf("%d 데미지 공격 준비 중", 10+floor),
		}
	} else {
		enemy.Intent = *intent
	}
	
	return enemy
}

// generateInitialIntent AI를 사용해서 초기 의도 생성
func (h *GameHandler) generateInitialIntent(enemy *domain.EnemyState, aiType string, floor int) (*domain.EnemyIntent, error) {
	// 기본 플레이어 상태 (초기 의도 계산용)
	dummyPlayer := &domain.PlayerState{
		Health:    80,
		MaxHealth: 80,
		Shield:    0,
		Energy:    3,
		ActivePowers: make(map[string]domain.PowerState),
		Buffs:        []domain.BuffState{},
		Debuffs:      []domain.DebuffState{},
	}
	
	// 기본 게임 상태
	dummyGameState := &domain.GameState{}
	
	return h.aiManager.CalculateNextIntent(
		enemy,
		dummyPlayer,
		dummyGameState,
		0, // 초기 턴
		floor,
		aiType,
	)
}

// getAITypeFromEnemyID 적 ID에서 AI 타입 추출
func (h *GameHandler) getAITypeFromEnemyID(enemyID string) string {
	// 적 ID에서 타입 추출 (예: "enemy_1_BRUTE" -> "aggressive")
	if contains(enemyID, "BRUTE") {
		return "aggressive"
	}
	if contains(enemyID, "GUARDIAN") {
		return "defensive"
	}
	if contains(enemyID, "ELITE") {
		return "balanced"
	}
	return "balanced" // 기본값
}

// contains 문자열 포함 여부 확인 헬퍼
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    str[:len(substr)] == substr || 
		    str[len(str)-len(substr):] == substr ||
		    indexOf(str, substr) >= 0)
}

// indexOf 문자열 인덱스 찾기 헬퍼
func indexOf(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func (h *GameHandler) generatePath(gameMode domain.GameMode) []domain.FloorNode {
	// TODO: Implement proper path generation
	return []domain.FloorNode{
		{
			ID:    "1-1",
			Type:  "COMBAT",
			Floor: 1,
			X:     0,
			Y:     0,
			NextNodes: []string{"2-1", "2-2"},
		},
	}
}

func (h *GameHandler) processPlayCard(session *domain.GameSession, playerState *domain.PlayerState, enemyState *domain.EnemyState, gameState *domain.GameState, cardID *string, targetID *string) (map[string]interface{}, error) {
	if cardID == nil {
		return nil, fmt.Errorf("카드 ID가 필요합니다")
	}

	// Check if card is in hand
	if !playerState.HasCardInHand(*cardID) {
		return nil, fmt.Errorf("손에 없는 카드입니다")
	}

	// Get card details
	card, err := h.cardRepo.GetByID(*cardID)
	if err != nil || card == nil {
		return nil, fmt.Errorf("카드 정보를 찾을 수 없습니다")
	}

	// Check if player has enough energy
	if !playerState.CanPlayCard(card) {
		return nil, fmt.Errorf("에너지가 부족합니다")
	}

	// Spend energy
	playerState.SpendEnergy(card.Cost)

	// Remove card from hand
	newHand := []string{}
	for _, id := range playerState.Hand {
		if id != *cardID {
			newHand = append(newHand, id)
		}
	}
	playerState.Hand = newHand

	// Process card effects using the effect executor
	executionResult, err := h.effectExecutor.ExecuteCardEffects(card, playerState, enemyState, gameState, targetID)
	if err != nil {
		return nil, fmt.Errorf("카드 효과 실행 실패: %w", err)
	}
	effects := executionResult.ToMap()

	// Add card to discard pile (unless it exhausts)
	if card.Type != domain.CardTypePower {
		playerState.DiscardPile = append(playerState.DiscardPile, *cardID)
	}

	// Update damage statistics
	if damage, ok := effects["damage_dealt"].(int); ok {
		session.DamageDealt += damage
	}

	return map[string]interface{}{
		"message": "카드를 사용했습니다",
		"card": card,
		"effects": effects,
		"energy_remaining": playerState.Energy,
	}, nil
}


func (h *GameHandler) processUsePotion(session *domain.GameSession, playerState *domain.PlayerState, enemyState *domain.EnemyState, gameState *domain.GameState, actionData json.RawMessage) (map[string]interface{}, error) {
	// TODO: Implement potion usage
	return nil, fmt.Errorf("포션 사용은 아직 구현되지 않았습니다")
}

func (h *GameHandler) processEnemyTurn(session *domain.GameSession, playerState *domain.PlayerState, enemyState *domain.EnemyState, gameState *domain.GameState) []map[string]interface{} {
	actions := []map[string]interface{}{}

	// AI 타입 결정 (적 ID에서 추출 또는 기본값)
	aiType := h.getAITypeFromEnemyID(enemyState.ID)
	
	// AI 시스템을 사용해서 적 턴 처리
	aiResult, err := h.aiManager.ProcessEnemyTurn(
		enemyState,
		playerState,
		gameState,
		session.CurrentTurn,
		session.CurrentFloor,
		aiType,
	)
	
	if err != nil {
		// AI 처리 실패시 기본 공격
		damage := 10 + session.CurrentFloor
		actualDamage := playerState.ApplyDamage(damage)
		session.DamageTaken += actualDamage
		
		actions = append(actions, map[string]interface{}{
			"type": "attack",
			"damage": damage,
			"actual_damage": actualDamage,
			"shield_blocked": damage - actualDamage,
			"message": fmt.Sprintf("적이 %d 데미지로 공격했습니다!", actualDamage),
		})
		
		// 기본 다음 의도 설정
		enemyState.Intent = domain.EnemyIntent{
			Type:        "ATTACK",
			Value:       damage,
			Description: fmt.Sprintf("%d 데미지 공격 준비 중", damage),
		}
	} else {
		// AI 결과를 액션으로 변환
		action := map[string]interface{}{
			"type": aiResult.Action.Type,
			"description": aiResult.Action.Description,
			"success": aiResult.Success,
		}
		
		// 결과에 따른 추가 정보
		if aiResult.Damage > 0 {
			action["damage"] = aiResult.Damage
			session.DamageTaken += aiResult.Damage
		}
		if aiResult.Shield > 0 {
			action["shield"] = aiResult.Shield
		}
		if len(aiResult.Buffs) > 0 {
			action["buffs_applied"] = aiResult.Buffs
		}
		if len(aiResult.Debuffs) > 0 {
			action["debuffs_applied"] = aiResult.Debuffs
		}
		if len(aiResult.Messages) > 0 {
			action["messages"] = aiResult.Messages
		}
		
		actions = append(actions, action)
		
		// 다음 의도 설정
		if aiResult.NextIntent != nil {
			enemyState.Intent = *aiResult.NextIntent
		}
	}

	return actions
}

func (h *GameHandler) processVictory(session *domain.GameSession, playerState *domain.PlayerState, enemyState *domain.EnemyState, gameState *domain.GameState) map[string]interface{} {
	// Calculate rewards
	goldReward := 50 + (session.CurrentFloor * 10)
	gameState.Gold += goldReward
	
	// Generate card rewards
	// TODO: Implement proper card reward generation
	cardRewards := []string{"card_004", "card_010", "card_015"}
	gameState.CardRewards = cardRewards

	// Update score
	session.Score += 100 + (session.CurrentFloor * 20)

	// Check if this was the boss
	if session.CurrentFloor%10 == 0 {
		// Game completed!
		session.Status = domain.GameStatusCompleted
		h.gameRepo.EndSession(session.ID, domain.GameStatusCompleted)
		h.userRepo.IncrementGamesWon(session.UserID)
		
		return map[string]interface{}{
			"message": "게임 클리어!",
			"result": "victory",
			"final_score": session.Score,
			"gold_reward": goldReward,
			"card_rewards": cardRewards,
		}
	}

	// Prepare for next floor
	session.CurrentFloor++
	gameState.FloorType = "REWARD"
	
	// Save state
	h.gameRepo.SaveGameState(session.ID, playerState, enemyState, gameState)
	h.gameRepo.UpdateSession(session)

	return map[string]interface{}{
		"message": "전투 승리!",
		"result": "floor_clear",
		"gold_reward": goldReward,
		"card_rewards": cardRewards,
		"next_floor": session.CurrentFloor,
		"player_state": playerState,
		"enemy_state": nil,
		"game_state": gameState,
	}
}

func (h *GameHandler) updateEffectDurations(playerState *domain.PlayerState, enemyState *domain.EnemyState) {
	// Update player buffs
	newBuffs := []domain.BuffState{}
	for _, buff := range playerState.Buffs {
		if buff.Duration > 0 {
			buff.Duration--
			if buff.Duration > 0 {
				newBuffs = append(newBuffs, buff)
			}
		} else if buff.Duration == -1 {
			newBuffs = append(newBuffs, buff)
		}
	}
	playerState.Buffs = newBuffs

	// Update player debuffs
	newDebuffs := []domain.DebuffState{}
	for _, debuff := range playerState.Debuffs {
		if debuff.Duration > 0 {
			debuff.Duration--
			if debuff.Duration > 0 {
				newDebuffs = append(newDebuffs, debuff)
			}
		} else if debuff.Duration == -1 {
			newDebuffs = append(newDebuffs, debuff)
		}
	}
	playerState.Debuffs = newDebuffs

	// Update enemy buffs and debuffs similarly
	// TODO: Implement enemy effect duration updates
}