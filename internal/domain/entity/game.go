package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GameStatus string
type GameMode string

const (
	GameStatusActive    GameStatus = "active"
	GameStatusCompleted GameStatus = "completed"
	GameStatusFailed    GameStatus = "failed"
	
	GameModeStory         GameMode = "story"
	GameModeDailyChallenge GameMode = "daily_challenge"
	GameModeEvent         GameMode = "event"
)

type GameSession struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Status    GameStatus      `json:"status" db:"status"`
	GameMode  GameMode        `json:"game_mode" db:"game_mode"`
	Floor     int             `json:"floor" db:"floor"`
	Score     int             `json:"score" db:"score"`
	GameState json.RawMessage `json:"game_state" db:"game_state"`
	StartedAt time.Time       `json:"started_at" db:"started_at"`
	EndedAt   *time.Time      `json:"ended_at,omitempty" db:"ended_at"`
}

type GameState struct {
	ID            string           `json:"id"`
	PlayerID      string           `json:"player_id"`
	Floor         int              `json:"floor"`
	Turn          int              `json:"turn"`
	Energy        Energy           `json:"energy"`
	Player        *PlayerState     `json:"player"`
	Enemies       []*EnemyState    `json:"enemies"`
	Hand          []*CardInstance  `json:"hand"`
	DrawPile      []*CardInstance  `json:"draw_pile"`
	DiscardPile   []*CardInstance  `json:"discard_pile"`
	ExhaustPile   []*CardInstance  `json:"exhaust_pile"`
	ActiveEffects []ActiveEffect   `json:"active_effects"`
	TurnHistory   []TurnAction     `json:"turn_history"`
}

type Energy struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

type PlayerState struct {
	HP      HealthPoint `json:"hp"`
	Shield  int         `json:"shield"`
	Buffs   []Buff      `json:"buffs"`
	Debuffs []Debuff    `json:"debuffs"`
}

type HealthPoint struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

type EnemyState struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	HP       HealthPoint `json:"hp"`
	Shield   int         `json:"shield"`
	Intent   EnemyIntent `json:"intent"`
	Buffs    []Buff      `json:"buffs"`
	Debuffs  []Debuff    `json:"debuffs"`
	Position int         `json:"position"`
}

type EnemyIntent struct {
	Type   IntentType `json:"type"`
	Value  int        `json:"value"`
	Target string     `json:"target,omitempty"`
}

type IntentType string

const (
	IntentAttack      IntentType = "ATTACK"
	IntentDefend      IntentType = "DEFEND"
	IntentBuff        IntentType = "BUFF"
	IntentDebuff      IntentType = "DEBUFF"
	IntentUnknown     IntentType = "UNKNOWN"
)

type Buff struct {
	Type     BuffType `json:"type"`
	Value    int      `json:"value"`
	Duration int      `json:"duration"`
}

type BuffType string

const (
	BuffStrength    BuffType = "STRENGTH"
	BuffDexterity   BuffType = "DEXTERITY"
	BuffRegeneration BuffType = "REGENERATION"
	BuffThorns      BuffType = "THORNS"
)

type Debuff struct {
	Type     DebuffType `json:"type"`
	Value    int        `json:"value"`
	Duration int        `json:"duration"`
}

type DebuffType string

const (
	DebuffVulnerable DebuffType = "VULNERABLE"
	DebuffWeak       DebuffType = "WEAK"
	DebuffPoison     DebuffType = "POISON"
	DebuffConfusion  DebuffType = "CONFUSION"
)

type ActiveEffect struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Source   string          `json:"source"`
	Target   string          `json:"target"`
	Value    int             `json:"value"`
	Duration int             `json:"duration"`
	Data     json.RawMessage `json:"data,omitempty"`
}

type TurnAction struct {
	Turn       int             `json:"turn"`
	Type       ActionType      `json:"type"`
	ActorID    string          `json:"actor_id"`
	TargetID   string          `json:"target_id,omitempty"`
	CardID     string          `json:"card_id,omitempty"`
	Result     json.RawMessage `json:"result"`
	Timestamp  time.Time       `json:"timestamp"`
}

type ActionType string

const (
	ActionPlayCard   ActionType = "PLAY_CARD"
	ActionEndTurn    ActionType = "END_TURN"
	ActionEnemyMove  ActionType = "ENEMY_MOVE"
	ActionTrigger    ActionType = "TRIGGER"
)

func NewGameSession(userID uuid.UUID, mode GameMode) *GameSession {
	return &GameSession{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    GameStatusActive,
		GameMode:  mode,
		Floor:     1,
		Score:     0,
		StartedAt: time.Now(),
	}
}

func (gs *GameSession) End(score int) {
	now := time.Now()
	gs.EndedAt = &now
	gs.Score = score
	if gs.Status == GameStatusActive {
		gs.Status = GameStatusCompleted
	}
}