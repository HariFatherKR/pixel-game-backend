package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/auth"
	"github.com/yourusername/pixel-game/internal/domain"
	"github.com/yourusername/pixel-game/internal/middleware"
)

// CardHandler handles card-related HTTP requests
type CardHandler struct {
	cardRepo   domain.CardRepository
	jwtManager *auth.JWTManager
}

// NewCardHandler creates a new card handler
func NewCardHandler(cardRepo domain.CardRepository, jwtManager *auth.JWTManager) *CardHandler {
	return &CardHandler{
		cardRepo:   cardRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers card routes
func (h *CardHandler) RegisterRoutes(router *gin.RouterGroup) {
	cards := router.Group("/cards")
	{
		// Public routes
		cards.GET("", h.GetCards)
		cards.GET("/:id", h.GetCard)
		
		// Protected routes
		protected := cards.Group("")
		protected.Use(middleware.AuthMiddleware(h.jwtManager))
		{
			protected.GET("/my-collection", h.GetMyCollection)
			protected.POST("/decks", h.CreateDeck)
			protected.GET("/decks", h.GetMyDecks)
			protected.GET("/decks/:id", h.GetDeck)
			protected.PUT("/decks/:id", h.UpdateDeck)
			protected.DELETE("/decks/:id", h.DeleteDeck)
			protected.PUT("/decks/:id/activate", h.ActivateDeck)
			protected.GET("/decks/active", h.GetActiveDeck)
		}
	}
}

// GetCards godoc
// @Summary 카드 목록 조회
// @Description 게임의 모든 카드 목록을 조회합니다. 필터와 페이지네이션을 지원합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Param type query string false "카드 타입 (ACTION, EVENT, POWER)"
// @Param rarity query string false "카드 희귀도 (COMMON, RARE, EPIC, LEGENDARY)"
// @Param min_cost query int false "최소 코스트"
// @Param max_cost query int false "최대 코스트"
// @Param search query string false "검색어 (카드 이름, 설명)"
// @Param limit query int false "결과 개수 제한" default(20)
// @Param offset query int false "결과 시작 위치" default(0)
// @Success 200 {object} map[string]interface{} "카드 목록"
// @Failure 400 {object} map[string]interface{} "잘못된 요청"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards [get]
func (h *CardHandler) GetCards(c *gin.Context) {
	filter := domain.CardFilter{
		Limit:  20,
		Offset: 0,
	}

	// Parse query parameters
	if cardType := c.Query("type"); cardType != "" {
		ct := domain.CardType(cardType)
		filter.Type = &ct
	}

	if rarity := c.Query("rarity"); rarity != "" {
		cr := domain.CardRarity(rarity)
		filter.Rarity = &cr
	}

	if minCost := c.Query("min_cost"); minCost != "" {
		if cost, err := strconv.Atoi(minCost); err == nil {
			filter.MinCost = &cost
		}
	}

	if maxCost := c.Query("max_cost"); maxCost != "" {
		if cost, err := strconv.Atoi(maxCost); err == nil {
			filter.MaxCost = &cost
		}
	}

	if search := c.Query("search"); search != "" {
		filter.SearchTerm = &search
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	cards, err := h.cardRepo.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "카드 목록을 조회할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards":  cards,
		"count":  len(cards),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetCard godoc
// @Summary 카드 상세 조회
// @Description 특정 카드의 상세 정보를 조회합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Param id path string true "카드 ID"
// @Success 200 {object} domain.Card "카드 정보"
// @Failure 404 {object} map[string]interface{} "카드를 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/{id} [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	cardID := c.Param("id")

	card, err := h.cardRepo.GetByID(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "카드 조회 중 오류가 발생했습니다",
		})
		return
	}

	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "카드를 찾을 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, card)
}

// GetMyCollection godoc
// @Summary 내 카드 컬렉션 조회
// @Description 현재 사용자가 보유한 카드 목록을 조회합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "카드 컬렉션"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/my-collection [get]
func (h *CardHandler) GetMyCollection(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	userCards, err := h.cardRepo.GetUserCards(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "카드 컬렉션을 조회할 수 없습니다",
		})
		return
	}

	// Group cards by type
	cardsByType := make(map[domain.CardType][]*domain.UserCard)
	for _, uc := range userCards {
		cardsByType[uc.Card.Type] = append(cardsByType[uc.Card.Type], uc)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       len(userCards),
		"cards":       userCards,
		"by_type":     cardsByType,
		"collection": gin.H{
			"action_cards": len(cardsByType[domain.CardTypeAction]),
			"event_cards":  len(cardsByType[domain.CardTypeEvent]),
			"power_cards":  len(cardsByType[domain.CardTypePower]),
		},
	})
}

// CreateDeck godoc
// @Summary 덱 생성
// @Description 새로운 덱을 생성합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param deck body map[string]interface{} true "덱 정보"
// @Success 201 {object} domain.Deck "생성된 덱"
// @Failure 400 {object} map[string]interface{} "잘못된 요청"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks [post]
func (h *CardHandler) CreateDeck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	var req struct {
		Name    string   `json:"name" binding:"required"`
		CardIDs []string `json:"card_ids" binding:"required,min=1,max=30"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청입니다",
		})
		return
	}

	// Validate deck size
	if len(req.CardIDs) < 10 || len(req.CardIDs) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "덱은 10장 이상 30장 이하의 카드로 구성되어야 합니다",
		})
		return
	}

	// Verify user owns all cards
	userCards, err := h.cardRepo.GetUserCards(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "카드 검증 중 오류가 발생했습니다",
		})
		return
	}

	userCardMap := make(map[string]bool)
	for _, uc := range userCards {
		userCardMap[uc.CardID] = true
	}

	for _, cardID := range req.CardIDs {
		if !userCardMap[cardID] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "보유하지 않은 카드가 포함되어 있습니다",
			})
			return
		}
	}

	deck := &domain.Deck{
		UserID:   userID.(int),
		Name:     req.Name,
		CardIDs:  req.CardIDs,
		IsActive: false,
	}

	if err := h.cardRepo.CreateDeck(deck); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 생성 중 오류가 발생했습니다",
		})
		return
	}

	c.JSON(http.StatusCreated, deck)
}

// GetMyDecks godoc
// @Summary 내 덱 목록 조회
// @Description 현재 사용자의 덱 목록을 조회합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "덱 목록"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks [get]
func (h *CardHandler) GetMyDecks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	decks, err := h.cardRepo.GetUserDecks(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 목록을 조회할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decks": decks,
		"count": len(decks),
	})
}

// GetDeck godoc
// @Summary 덱 상세 조회
// @Description 특정 덱의 상세 정보를 조회합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "덱 ID"
// @Success 200 {object} map[string]interface{} "덱 정보"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "덱을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks/{id} [get]
func (h *CardHandler) GetDeck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	deckID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 덱 ID입니다",
		})
		return
	}

	deck, err := h.cardRepo.GetDeck(deckID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 조회 중 오류가 발생했습니다",
		})
		return
	}

	if deck == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "덱을 찾을 수 없습니다",
		})
		return
	}

	if deck.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 덱에 접근할 권한이 없습니다",
		})
		return
	}

	// Get card details
	cards, err := h.cardRepo.GetByIDs(deck.CardIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "카드 정보를 조회할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deck":  deck,
		"cards": cards,
	})
}

// UpdateDeck godoc
// @Summary 덱 수정
// @Description 덱 정보를 수정합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "덱 ID"
// @Param deck body map[string]interface{} true "수정할 덱 정보"
// @Success 200 {object} domain.Deck "수정된 덱"
// @Failure 400 {object} map[string]interface{} "잘못된 요청"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "덱을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks/{id} [put]
func (h *CardHandler) UpdateDeck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	deckID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 덱 ID입니다",
		})
		return
	}

	var req struct {
		Name    string   `json:"name"`
		CardIDs []string `json:"card_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청입니다",
		})
		return
	}

	deck, err := h.cardRepo.GetDeck(deckID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 조회 중 오류가 발생했습니다",
		})
		return
	}

	if deck == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "덱을 찾을 수 없습니다",
		})
		return
	}

	if deck.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 덱을 수정할 권한이 없습니다",
		})
		return
	}

	// Update fields
	if req.Name != "" {
		deck.Name = req.Name
	}

	if len(req.CardIDs) > 0 {
		if len(req.CardIDs) < 10 || len(req.CardIDs) > 30 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "덱은 10장 이상 30장 이하의 카드로 구성되어야 합니다",
			})
			return
		}

		// Verify user owns all cards
		userCards, err := h.cardRepo.GetUserCards(userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "카드 검증 중 오류가 발생했습니다",
			})
			return
		}

		userCardMap := make(map[string]bool)
		for _, uc := range userCards {
			userCardMap[uc.CardID] = true
		}

		for _, cardID := range req.CardIDs {
			if !userCardMap[cardID] {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "보유하지 않은 카드가 포함되어 있습니다",
				})
				return
			}
		}

		deck.CardIDs = req.CardIDs
	}

	if err := h.cardRepo.UpdateDeck(deck); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 수정 중 오류가 발생했습니다",
		})
		return
	}

	c.JSON(http.StatusOK, deck)
}

// DeleteDeck godoc
// @Summary 덱 삭제
// @Description 덱을 삭제합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "덱 ID"
// @Success 204 "삭제 성공"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "덱을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks/{id} [delete]
func (h *CardHandler) DeleteDeck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	deckID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 덱 ID입니다",
		})
		return
	}

	deck, err := h.cardRepo.GetDeck(deckID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 조회 중 오류가 발생했습니다",
		})
		return
	}

	if deck == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "덱을 찾을 수 없습니다",
		})
		return
	}

	if deck.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 덱을 삭제할 권한이 없습니다",
		})
		return
	}

	if err := h.cardRepo.DeleteDeck(deckID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 삭제 중 오류가 발생했습니다",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ActivateDeck godoc
// @Summary 덱 활성화
// @Description 덱을 활성 덱으로 설정합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "덱 ID"
// @Success 200 {object} map[string]interface{} "성공 메시지"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 403 {object} map[string]interface{} "권한 없음"
// @Failure 404 {object} map[string]interface{} "덱을 찾을 수 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks/{id}/activate [put]
func (h *CardHandler) ActivateDeck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	deckID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 덱 ID입니다",
		})
		return
	}

	deck, err := h.cardRepo.GetDeck(deckID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 조회 중 오류가 발생했습니다",
		})
		return
	}

	if deck == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "덱을 찾을 수 없습니다",
		})
		return
	}

	if deck.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "이 덱을 활성화할 권한이 없습니다",
		})
		return
	}

	if err := h.cardRepo.SetActiveDeck(userID.(int), deckID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "덱 활성화 중 오류가 발생했습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "덱이 활성화되었습니다",
		"deck_id": deckID,
	})
}

// GetActiveDeck godoc
// @Summary 활성 덱 조회
// @Description 현재 활성화된 덱을 조회합니다.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "활성 덱 정보"
// @Failure 401 {object} map[string]interface{} "인증 필요"
// @Failure 404 {object} map[string]interface{} "활성 덱이 없음"
// @Failure 500 {object} map[string]interface{} "서버 에러"
// @Router /api/v1/cards/decks/active [get]
func (h *CardHandler) GetActiveDeck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "인증이 필요합니다",
		})
		return
	}

	deck, err := h.cardRepo.GetActiveDeck(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "활성 덱 조회 중 오류가 발생했습니다",
		})
		return
	}

	if deck == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "활성화된 덱이 없습니다",
		})
		return
	}

	// Get card details
	cards, err := h.cardRepo.GetByIDs(deck.CardIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "카드 정보를 조회할 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deck":  deck,
		"cards": cards,
	})
}