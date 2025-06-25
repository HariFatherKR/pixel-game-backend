package effects

import (
	"github.com/yourusername/pixel-game/internal/domain"
)

// EffectContext provides the context for executing card effects
type EffectContext struct {
	Session     *domain.GameSession
	PlayerState *domain.PlayerState
	EnemyState  *domain.EnemyState
	GameState   *domain.GameState
	SourceCard  *domain.Card
	TargetID    string // Could be enemy ID, card ID, etc.
}

// EffectResult contains the result of executing an effect
type EffectResult struct {
	Success      bool
	Damage       int
	Healing      int
	ShieldGained int
	CardsDrawn   []string
	BuffsApplied []domain.BuffState
	DebuffsApplied []domain.DebuffState
	EnergyUsed   int
	Messages     []string
	Errors       []string
}

// CardEffect is the interface that all card effects must implement
type CardEffect interface {
	// Execute runs the card effect
	Execute(ctx *EffectContext) (*EffectResult, error)
	
	// CanExecute checks if the effect can be executed in the current context
	CanExecute(ctx *EffectContext) (bool, string)
	
	// GetType returns the effect type
	GetType() string
	
	// GetDescription returns a description of what the effect does
	GetDescription() string
}

