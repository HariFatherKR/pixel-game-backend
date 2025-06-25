package entity

import (
	"encoding/json"
	"time"
)

type CardType string
type CardRarity string

const (
	CardTypeAction CardType = "ACTION"
	CardTypeEvent  CardType = "EVENT"
	
	RarityCommon    CardRarity = "COMMON"
	RarityRare      CardRarity = "RARE"
	RarityEpic      CardRarity = "EPIC"
	RarityLegendary CardRarity = "LEGENDARY"
)

type Card struct {
	ID            string          `json:"id" db:"id"`
	Name          string          `json:"name" db:"name"`
	Type          CardType        `json:"type" db:"type"`
	Rarity        CardRarity      `json:"rarity" db:"rarity"`
	Cost          int             `json:"cost" db:"cost"`
	Description   string          `json:"description" db:"description"`
	CodeSnippet   string          `json:"code_snippet" db:"code_snippet"`
	Effects       []Effect        `json:"effects" db:"effects"`
	VisualEffects json.RawMessage `json:"visual_effects" db:"visual_effects"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

type Effect struct {
	Type       EffectType             `json:"type"`
	Target     TargetType             `json:"target"`
	Value      int                    `json:"value"`
	Duration   int                    `json:"duration,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
}

type EffectType string

const (
	EffectDamage      EffectType = "DAMAGE"
	EffectShield      EffectType = "SHIELD"
	EffectHeal        EffectType = "HEAL"
	EffectDraw        EffectType = "DRAW"
	EffectEnergy      EffectType = "ENERGY"
	EffectBuff        EffectType = "BUFF"
	EffectDebuff      EffectType = "DEBUFF"
	EffectDiscard     EffectType = "DISCARD"
	EffectExhaust     EffectType = "EXHAUST"
)

type TargetType string

const (
	TargetSelf       TargetType = "SELF"
	TargetEnemy      TargetType = "ENEMY"
	TargetAllEnemies TargetType = "ALL_ENEMIES"
	TargetRandom     TargetType = "RANDOM_ENEMY"
)

type UserCard struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	CardID       string    `json:"card_id" db:"card_id"`
	AcquiredAt   time.Time `json:"acquired_at" db:"acquired_at"`
	IsUpgraded   bool      `json:"is_upgraded" db:"is_upgraded"`
	UpgradePath  string    `json:"upgrade_path,omitempty" db:"upgrade_path"`
	Level        int       `json:"level" db:"level"`
}

type CardInstance struct {
	Card       *Card  `json:"card"`
	InstanceID string `json:"instance_id"`
	IsUpgraded bool   `json:"is_upgraded"`
	TempCost   *int   `json:"temp_cost,omitempty"` // Temporary cost modification
}

func NewCardInstance(card *Card, upgraded bool) *CardInstance {
	return &CardInstance{
		Card:       card,
		InstanceID: generateInstanceID(),
		IsUpgraded: upgraded,
	}
}

func generateInstanceID() string {
	// Simple implementation - in production use UUID
	return time.Now().Format("20060102150405") + "-" + string(time.Now().UnixNano()%1000)
}