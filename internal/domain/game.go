package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Game status
type GameStatus string

const (
	GameStatusActive    GameStatus = "ACTIVE"
	GameStatusCompleted GameStatus = "COMPLETED"
	GameStatusFailed    GameStatus = "FAILED"
	GameStatusPaused    GameStatus = "PAUSED"
)

// Game mode
type GameMode string

const (
	GameModeStory         GameMode = "STORY"
	GameModeDailyChallenge GameMode = "DAILY_CHALLENGE"
	GameModeEvent         GameMode = "EVENT"
)

// Turn phase
type TurnPhase string

const (
	TurnPhaseStart  TurnPhase = "START"
	TurnPhaseDraw   TurnPhase = "DRAW"
	TurnPhaseMain   TurnPhase = "MAIN"
	TurnPhaseEnd    TurnPhase = "END"
	TurnPhaseEnemy  TurnPhase = "ENEMY"
)

// GameSession represents an active game session
type GameSession struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          int             `json:"user_id" db:"user_id"`
	Status          GameStatus      `json:"status" db:"status"`
	GameMode        GameMode        `json:"game_mode" db:"game_mode"`
	CurrentFloor    int             `json:"current_floor" db:"current_floor"`
	CurrentTurn     int             `json:"current_turn" db:"current_turn"`
	TurnPhase       TurnPhase       `json:"turn_phase" db:"turn_phase"`
	PlayerState     json.RawMessage `json:"player_state" db:"player_state"`
	EnemyState      json.RawMessage `json:"enemy_state" db:"enemy_state"`
	GameState       json.RawMessage `json:"game_state" db:"game_state"`
	DeckSnapshot    []string        `json:"deck_snapshot" db:"deck_snapshot"`
	Score           int             `json:"score" db:"score"`
	CardsPlayed     int             `json:"cards_played" db:"cards_played"`
	DamageDealt     int             `json:"damage_dealt" db:"damage_dealt"`
	DamageTaken     int             `json:"damage_taken" db:"damage_taken"`
	StartedAt       time.Time       `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time      `json:"completed_at" db:"completed_at"`
	LastActionAt    time.Time       `json:"last_action_at" db:"last_action_at"`
	TurnTimeLimit   int             `json:"turn_time_limit" db:"turn_time_limit"` // seconds
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// PlayerState represents the player's current state in a game
type PlayerState struct {
	Health       int                   `json:"health"`
	MaxHealth    int                   `json:"max_health"`
	Shield       int                   `json:"shield"`
	Energy       int                   `json:"energy"`
	MaxEnergy    int                   `json:"max_energy"`
	Hand         []string              `json:"hand"`         // Card IDs in hand
	DrawPile     []string              `json:"draw_pile"`    // Card IDs in draw pile
	DiscardPile  []string              `json:"discard_pile"` // Card IDs in discard pile
	ExhaustPile  []string              `json:"exhaust_pile"` // Card IDs removed from play
	Deck         []string              `json:"deck"`         // Player's deck card IDs
	ActivePowers map[string]PowerState `json:"active_powers"`
	Buffs        []BuffState           `json:"buffs"`
	Debuffs      []DebuffState         `json:"debuffs"`
}

// EnemyState represents an enemy's current state
type EnemyState struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Health       int           `json:"health"`
	MaxHealth    int           `json:"max_health"`
	Shield       int           `json:"shield"`
	Intent       EnemyIntent   `json:"intent"`
	ActivePowers []PowerState  `json:"active_powers"`
	Buffs        []BuffState   `json:"buffs"`
	Debuffs      []DebuffState `json:"debuffs"`
}

// EnemyIntent represents what the enemy plans to do
type EnemyIntent struct {
	Type        string `json:"type"` // ATTACK, DEFEND, BUFF, DEBUFF, UNKNOWN
	Value       int    `json:"value"`
	Description string `json:"description"`
}

// PowerState represents an active power effect
type PowerState struct {
	PowerID     string `json:"power_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stacks      int    `json:"stacks"`
	Duration    int    `json:"duration"` // -1 for permanent
}

// BuffState represents a temporary positive effect
type BuffState struct {
	BuffID      string `json:"buff_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Duration    int    `json:"duration"`
}

// DebuffState represents a temporary negative effect
type DebuffState struct {
	DebuffID    string `json:"debuff_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Duration    int    `json:"duration"`
}

// GameState contains overall game progression data
type GameState struct {
	FloorType     string                 `json:"floor_type"` // COMBAT, EVENT, SHOP, REST, BOSS
	FloorData     map[string]interface{} `json:"floor_data"`
	Relics        []string               `json:"relics"`
	Gold          int                    `json:"gold"`
	PotionSlots   int                    `json:"potion_slots"`
	Potions       []string               `json:"potions"`
	CardRewards   []string               `json:"card_rewards"`
	Path          []FloorNode            `json:"path"`
	CurrentNodeID string                 `json:"current_node_id"`
}

// FloorNode represents a node in the game map
type FloorNode struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Floor    int    `json:"floor"`
	Visited  bool   `json:"visited"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	NextNodes []string `json:"next_nodes"`
}

// GameAction represents a player action in the game
type GameAction struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	SessionID  uuid.UUID       `json:"session_id" db:"session_id"`
	ActionType string          `json:"action_type" db:"action_type"`
	CardID     *string         `json:"card_id,omitempty" db:"card_id"`
	TargetID   *string         `json:"target_id,omitempty" db:"target_id"`
	ActionData json.RawMessage `json:"action_data" db:"action_data"`
	Timestamp  time.Time       `json:"timestamp" db:"timestamp"`
}

// ActionType represents different types of game actions
type ActionType string

const (
	ActionTypePlayCard   ActionType = "PLAY_CARD"
	ActionTypeEndTurn    ActionType = "END_TURN"
	ActionTypeUsePotion  ActionType = "USE_POTION"
	ActionTypeSelectCard ActionType = "SELECT_CARD"
	ActionTypeSelectPath ActionType = "SELECT_PATH"
	ActionTypeRest       ActionType = "REST"
	ActionTypeShop       ActionType = "SHOP"
	ActionTypeSkip       ActionType = "SKIP"
)

// GameRepository interface
type GameRepository interface {
	// Session management
	CreateSession(session *GameSession) error
	GetSession(sessionID uuid.UUID) (*GameSession, error)
	GetActiveSession(userID int) (*GameSession, error)
	UpdateSession(session *GameSession) error
	EndSession(sessionID uuid.UUID, status GameStatus) error
	
	// Game state
	SaveGameState(sessionID uuid.UUID, playerState *PlayerState, enemyState *EnemyState, gameState *GameState) error
	LoadGameState(sessionID uuid.UUID) (*PlayerState, *EnemyState, *GameState, error)
	
	// Actions
	RecordAction(action *GameAction) error
	GetSessionActions(sessionID uuid.UUID) ([]*GameAction, error)
	
	// Statistics
	GetUserGameStats(userID int) (*UserGameStats, error)
	UpdateGameStats(sessionID uuid.UUID) error
}

// UserGameStats represents aggregated game statistics for a user
type UserGameStats struct {
	TotalGames      int     `json:"total_games"`
	GamesWon        int     `json:"games_won"`
	GamesLost       int     `json:"games_lost"`
	WinRate         float64 `json:"win_rate"`
	HighestFloor    int     `json:"highest_floor"`
	TotalScore      int     `json:"total_score"`
	HighestScore    int     `json:"highest_score"`
	TotalPlayTime   int     `json:"total_play_time"` // seconds
	AverageGameTime int     `json:"average_game_time"`
	FavoriteCards   []string `json:"favorite_cards"`
}

// Helper methods
func (gs *GameSession) IsActive() bool {
	return gs.Status == GameStatusActive
}

func (gs *GameSession) CanTakeAction() bool {
	return gs.Status == GameStatusActive && 
		(gs.TurnPhase == TurnPhaseMain || gs.TurnPhase == TurnPhaseDraw)
}

func (ps *PlayerState) CanPlayCard(card *Card) bool {
	return ps.Energy >= card.Cost
}

func (ps *PlayerState) HasCardInHand(cardID string) bool {
	for _, id := range ps.Hand {
		if id == cardID {
			return true
		}
	}
	return false
}

func (ps *PlayerState) ApplyDamage(damage int) int {
	// Apply damage to shield first
	if ps.Shield > 0 {
		if ps.Shield >= damage {
			ps.Shield -= damage
			return 0
		}
		damage -= ps.Shield
		ps.Shield = 0
	}
	
	// Apply remaining damage to health
	ps.Health -= damage
	if ps.Health < 0 {
		ps.Health = 0
	}
	
	return damage
}

func (ps *PlayerState) Heal(amount int) {
	ps.Health += amount
	if ps.Health > ps.MaxHealth {
		ps.Health = ps.MaxHealth
	}
}

func (ps *PlayerState) GainShield(amount int) {
	ps.Shield += amount
}

func (ps *PlayerState) SpendEnergy(amount int) bool {
	if ps.Energy >= amount {
		ps.Energy -= amount
		return true
	}
	return false
}

func (ps *PlayerState) DrawCards(count int) []string {
	drawn := []string{}
	
	for i := 0; i < count && len(ps.Hand) < 10; i++ {
		if len(ps.DrawPile) == 0 {
			// Shuffle discard pile into draw pile
			ps.DrawPile = ps.DiscardPile
			ps.DiscardPile = []string{}
			// TODO: Implement shuffle
		}
		
		if len(ps.DrawPile) > 0 {
			card := ps.DrawPile[0]
			ps.DrawPile = ps.DrawPile[1:]
			ps.Hand = append(ps.Hand, card)
			drawn = append(drawn, card)
		}
	}
	
	return drawn
}