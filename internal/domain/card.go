package domain

import (
	"encoding/json"
	"time"
)

// Card types
type CardType string
type CardRarity string

const (
	CardTypeAction CardType = "ACTION"
	CardTypeEvent  CardType = "EVENT"
	CardTypePower  CardType = "POWER"
)

const (
	CardRarityCommon    CardRarity = "COMMON"
	CardRarityRare      CardRarity = "RARE"
	CardRarityEpic      CardRarity = "EPIC"
	CardRarityLegendary CardRarity = "LEGENDARY"
)

// Card represents a game card with code-based effects
type Card struct {
	ID            string          `json:"id" db:"id"`
	Name          string          `json:"name" db:"name"`
	Type          CardType        `json:"type" db:"type"`
	Rarity        CardRarity      `json:"rarity" db:"rarity"`
	Cost          int             `json:"cost" db:"cost"`
	Description   string          `json:"description" db:"description"`
	CodeSnippet   string          `json:"code_snippet" db:"code_snippet"`
	Effects       json.RawMessage `json:"effects" db:"effects"`
	VisualEffects json.RawMessage `json:"visual_effects" db:"visual_effects"`
	ImageURL      string          `json:"image_url" db:"image_url"`
	BaseDamage    int             `json:"base_damage" db:"base_damage"`
	BaseBlock     int             `json:"base_block" db:"base_block"`
	DrawAmount    int             `json:"draw_amount" db:"draw_amount"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

// UserCard represents a card owned by a user
type UserCard struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CardID      string    `json:"card_id" db:"card_id"`
	AcquiredAt  time.Time `json:"acquired_at" db:"acquired_at"`
	IsUpgraded  bool      `json:"is_upgraded" db:"is_upgraded"`
	UpgradePath string    `json:"upgrade_path" db:"upgrade_path"`
	Level       int       `json:"level" db:"level"`
	Card        *Card     `json:"card,omitempty"`
}

// CardEffect represents the structured effect data
type CardEffect struct {
	Type       string          `json:"type"`
	Target     string          `json:"target"`
	Value      int             `json:"value"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
	Conditions json.RawMessage `json:"conditions,omitempty"`
}

// VisualEffect represents the visual effect data
type VisualEffect struct {
	Action    string                 `json:"action"`
	Element   string                 `json:"element,omitempty"`
	Selector  string                 `json:"selector,omitempty"`
	Class     string                 `json:"class,omitempty"`
	Target    string                 `json:"target,omitempty"`
	Animation string                 `json:"animation,omitempty"`
	Style     map[string]interface{} `json:"style,omitempty"`
	Duration  int                    `json:"duration,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

// Deck represents a user's card deck
type Deck struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	CardIDs   []string  `json:"card_ids" db:"card_ids"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CardFilter for querying cards
type CardFilter struct {
	Type       *CardType
	Rarity     *CardRarity
	MaxCost    *int
	MinCost    *int
	SearchTerm *string
	Limit      int
	Offset     int
}

// CardRepository interface
type CardRepository interface {
	// Card master data operations
	GetAll(filter CardFilter) ([]*Card, error)
	GetByID(id string) (*Card, error)
	GetByIDs(ids []string) ([]*Card, error)
	Create(card *Card) error
	Update(card *Card) error
	Delete(id string) error
	
	// User card collection operations
	GetUserCards(userID int) ([]*UserCard, error)
	GetUserCard(userID int, cardID string) (*UserCard, error)
	AddCardToUser(userCard *UserCard) error
	UpdateUserCard(userCard *UserCard) error
	RemoveCardFromUser(userID int, cardID string) error
	
	// Deck operations
	CreateDeck(deck *Deck) error
	GetUserDecks(userID int) ([]*Deck, error)
	GetDeck(deckID int) (*Deck, error)
	UpdateDeck(deck *Deck) error
	DeleteDeck(deckID int) error
	SetActiveDeck(userID int, deckID int) error
	GetActiveDeck(userID int) (*Deck, error)
}

// Helper methods
func (c *Card) GetEffects() ([]CardEffect, error) {
	var effects []CardEffect
	if err := json.Unmarshal(c.Effects, &effects); err != nil {
		return nil, err
	}
	return effects, nil
}

func (c *Card) GetVisualEffects() (*VisualEffect, error) {
	var effect VisualEffect
	if c.VisualEffects == nil {
		return nil, nil
	}
	if err := json.Unmarshal(c.VisualEffects, &effect); err != nil {
		return nil, err
	}
	return &effect, nil
}

func (c *Card) IsPlayable(currentEnergy int) bool {
	return currentEnergy >= c.Cost
}

func (c *Card) GetRarityColor() string {
	switch c.Rarity {
	case CardRarityCommon:
		return "#808080" // Gray
	case CardRarityRare:
		return "#0080FF" // Blue
	case CardRarityEpic:
		return "#A335EE" // Purple
	case CardRarityLegendary:
		return "#FF8000" // Orange
	default:
		return "#FFFFFF"
	}
}