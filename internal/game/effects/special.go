package effects

import (
	"fmt"
	"math/rand"
	"github.com/yourusername/pixel-game/internal/domain"
)

// EnergyGainEffect gives additional energy
type EnergyGainEffect struct {
	amount int
}

// NewEnergyGainEffect creates an energy gain effect
func NewEnergyGainEffect(amount int) *EnergyGainEffect {
	return &EnergyGainEffect{amount: amount}
}

// Execute grants energy
func (e *EnergyGainEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	ctx.PlayerState.Energy += e.amount
	result.Messages = append(result.Messages, 
		fmt.Sprintf("Gained %d energy", e.amount))

	return result, nil
}

// CanExecute checks if energy can be gained
func (e *EnergyGainEffect) CanExecute(ctx *EffectContext) (bool, string) {
	return true, ""
}

// GetType returns the effect type
func (e *EnergyGainEffect) GetType() string {
	return "energy_gain"
}

// GetDescription returns the effect description
func (e *EnergyGainEffect) GetDescription() string {
	return fmt.Sprintf("Gain %d energy", e.amount)
}

// HealEffect restores health
type HealEffect struct {
	amount int
}

// NewHealEffect creates a heal effect
func NewHealEffect(amount int) *HealEffect {
	return &HealEffect{amount: amount}
}

// Execute heals the player
func (e *HealEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Calculate actual healing
	healAmount := e.amount
	if ctx.PlayerState.Health + healAmount > ctx.PlayerState.MaxHealth {
		healAmount = ctx.PlayerState.MaxHealth - ctx.PlayerState.Health
	}

	ctx.PlayerState.Health += healAmount
	result.Healing = healAmount
	result.Messages = append(result.Messages, 
		fmt.Sprintf("Healed %d health", healAmount))

	return result, nil
}

// CanExecute checks if healing can be done
func (e *HealEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if ctx.PlayerState.Health >= ctx.PlayerState.MaxHealth {
		return false, "already at full health"
	}
	return true, ""
}

// GetType returns the effect type
func (e *HealEffect) GetType() string {
	return "heal"
}

// GetDescription returns the effect description
func (e *HealEffect) GetDescription() string {
	return fmt.Sprintf("Heal %d health", e.amount)
}

// ExhaustEffect exhausts cards from hand
type ExhaustEffect struct {
	count     int
	targetSelf bool // If true, exhausts the played card itself
}

// NewExhaustEffect creates an exhaust effect
func NewExhaustEffect(count int, targetSelf bool) *ExhaustEffect {
	return &ExhaustEffect{
		count:      count,
		targetSelf: targetSelf,
	}
}

// Execute exhausts cards
func (e *ExhaustEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	if e.targetSelf && ctx.SourceCard != nil {
		// Exhaust the played card
		// This would be handled by the game logic after effect execution
		result.Messages = append(result.Messages, "This card will be exhausted")
	} else {
		// Exhaust random cards from hand
		exhausted := 0
		for i := 0; i < e.count && len(ctx.PlayerState.Hand) > 0; i++ {
			// Remove random card from hand
			randIndex := rand.Intn(len(ctx.PlayerState.Hand))
			card := ctx.PlayerState.Hand[randIndex]
			ctx.PlayerState.Hand = append(
				ctx.PlayerState.Hand[:randIndex],
				ctx.PlayerState.Hand[randIndex+1:]...)
			ctx.PlayerState.ExhaustPile = append(ctx.PlayerState.ExhaustPile, card)
			exhausted++
		}
		
		if exhausted > 0 {
			result.Messages = append(result.Messages, 
				fmt.Sprintf("Exhausted %d cards", exhausted))
		}
	}

	return result, nil
}

// CanExecute checks if cards can be exhausted
func (e *ExhaustEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if !e.targetSelf && len(ctx.PlayerState.Hand) == 0 {
		return false, "no cards in hand to exhaust"
	}
	return true, ""
}

// GetType returns the effect type
func (e *ExhaustEffect) GetType() string {
	return "exhaust"
}

// GetDescription returns the effect description
func (e *ExhaustEffect) GetDescription() string {
	if e.targetSelf {
		return "Exhaust this card"
	}
	return fmt.Sprintf("Exhaust %d cards from hand", e.count)
}

// RetainEffect allows cards to be retained between turns
type RetainEffect struct {
	cardID string
}

// NewRetainEffect creates a retain effect
func NewRetainEffect(cardID string) *RetainEffect {
	return &RetainEffect{cardID: cardID}
}

// Execute marks a card for retention
func (e *RetainEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Add retain buff to the card
	buff := domain.BuffState{
		BuffID:      fmt.Sprintf("retain_%s", e.cardID),
		Name:        "Retain",
		Description: "This card will not be discarded at end of turn",
		Value:       1,
		Duration:    1,
	}

	ctx.PlayerState.Buffs = append(ctx.PlayerState.Buffs, buff)
	result.BuffsApplied = append(result.BuffsApplied, buff)
	result.Messages = append(result.Messages, "This card will be retained")

	return result, nil
}

// CanExecute checks if retain can be applied
func (e *RetainEffect) CanExecute(ctx *EffectContext) (bool, string) {
	return true, ""
}

// GetType returns the effect type
func (e *RetainEffect) GetType() string {
	return "retain"
}

// GetDescription returns the effect description
func (e *RetainEffect) GetDescription() string {
	return "Retain this card"
}

// DoublePlayEffect plays the next card twice
type DoublePlayEffect struct{}

// NewDoublePlayEffect creates a double play effect
func NewDoublePlayEffect() *DoublePlayEffect {
	return &DoublePlayEffect{}
}

// Execute applies double play buff
func (e *DoublePlayEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Add double play buff
	buff := domain.BuffState{
		BuffID:      "double_play",
		Name:        "Double Play",
		Description: "Next card is played twice",
		Value:       1,
		Duration:    1, // Until next card is played
	}

	ctx.PlayerState.Buffs = append(ctx.PlayerState.Buffs, buff)
	result.BuffsApplied = append(result.BuffsApplied, buff)
	result.Messages = append(result.Messages, "Next card will be played twice")

	return result, nil
}

// CanExecute checks if double play can be applied
func (e *DoublePlayEffect) CanExecute(ctx *EffectContext) (bool, string) {
	// Check if already has double play
	for _, buff := range ctx.PlayerState.Buffs {
		if buff.BuffID == "double_play" {
			return false, "double play already active"
		}
	}
	return true, ""
}

// GetType returns the effect type
func (e *DoublePlayEffect) GetType() string {
	return "double_play"
}

// GetDescription returns the effect description
func (e *DoublePlayEffect) GetDescription() string {
	return "Next card is played twice"
}