package effects

import (
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// ShieldEffect implements shield/block gaining
type ShieldEffect struct {
	baseShield int
}

// NewShieldEffect creates a new shield effect
func NewShieldEffect(shield int) *ShieldEffect {
	return &ShieldEffect{
		baseShield: shield,
	}
}

// Execute grants shield to the player
func (e *ShieldEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Calculate actual shield (can be modified by buffs)
	shield := e.calculateShield(ctx)

	// Apply shield to player
	ctx.PlayerState.Shield += shield
	result.ShieldGained = shield
	result.Messages = append(result.Messages, 
		fmt.Sprintf("Gained %d shield", shield))

	return result, nil
}

// CanExecute checks if shield can be gained
func (e *ShieldEffect) CanExecute(ctx *EffectContext) (bool, string) {
	// Shield can always be gained
	return true, ""
}

// GetType returns the effect type
func (e *ShieldEffect) GetType() string {
	return "shield"
}

// GetDescription returns the effect description
func (e *ShieldEffect) GetDescription() string {
	return fmt.Sprintf("Gain %d shield", e.baseShield)
}

// calculateShield calculates actual shield considering modifiers
func (e *ShieldEffect) calculateShield(ctx *EffectContext) int {
	shield := e.baseShield

	// Check for dexterity buff
	if dexterity, exists := ctx.PlayerState.ActivePowers["dexterity"]; exists {
		shield += dexterity.Stacks
	}

	// Check for frail debuff
	for _, debuff := range ctx.PlayerState.Debuffs {
		if debuff.DebuffID == "frail" {
			shield = int(float64(shield) * 0.75)
			break
		}
	}

	return shield
}

// ReflectShieldEffect implements shield that reflects damage
type ReflectShieldEffect struct {
	baseShield int
	reflectPercent float64
}

// NewReflectShieldEffect creates a reflective shield effect
func NewReflectShieldEffect(shield int, reflectPercent float64) *ReflectShieldEffect {
	return &ReflectShieldEffect{
		baseShield: shield,
		reflectPercent: reflectPercent,
	}
}

// Execute grants reflective shield
func (e *ReflectShieldEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	// First apply normal shield
	shieldResult, err := NewShieldEffect(e.baseShield).Execute(ctx)
	if err != nil {
		return shieldResult, err
	}

	// Add reflection buff
	reflectBuff := domain.BuffState{
		BuffID:      "thorns",
		Name:        "Thorns",
		Description: fmt.Sprintf("Reflects %d%% of received damage", int(e.reflectPercent*100)),
		Value:       int(e.reflectPercent * 100),
		Duration:    1, // Lasts until next turn
	}

	ctx.PlayerState.Buffs = append(ctx.PlayerState.Buffs, reflectBuff)
	shieldResult.BuffsApplied = append(shieldResult.BuffsApplied, reflectBuff)
	shieldResult.Messages = append(shieldResult.Messages, 
		fmt.Sprintf("Gained thorns effect (%d%% reflection)", int(e.reflectPercent*100)))

	return shieldResult, nil
}

// CanExecute checks if reflective shield can be gained
func (e *ReflectShieldEffect) CanExecute(ctx *EffectContext) (bool, string) {
	return true, ""
}

// GetType returns the effect type
func (e *ReflectShieldEffect) GetType() string {
	return "reflect_shield"
}

// GetDescription returns the effect description
func (e *ReflectShieldEffect) GetDescription() string {
	return fmt.Sprintf("Gain %d shield and %d%% damage reflection", 
		e.baseShield, int(e.reflectPercent*100))
}

// BarricadeEffect makes shield not expire at end of turn
type BarricadeEffect struct{}

// NewBarricadeEffect creates a barricade effect
func NewBarricadeEffect() *BarricadeEffect {
	return &BarricadeEffect{}
}

// Execute applies barricade power
func (e *BarricadeEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Add barricade power
	if _, exists := ctx.PlayerState.ActivePowers["barricade"]; !exists {
		ctx.PlayerState.ActivePowers["barricade"] = domain.PowerState{
			PowerID:     "barricade",
			Name:        "Barricade",
			Description: "Shield no longer expires at the end of turn",
			Stacks:      1,
			Duration:    -1, // Permanent
		}
		result.Messages = append(result.Messages, "Shield will no longer expire at end of turn")
	}

	return result, nil
}

// CanExecute checks if barricade can be applied
func (e *BarricadeEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if _, exists := ctx.PlayerState.ActivePowers["barricade"]; exists {
		return false, "barricade already active"
	}
	return true, ""
}

// GetType returns the effect type
func (e *BarricadeEffect) GetType() string {
	return "barricade"
}

// GetDescription returns the effect description
func (e *BarricadeEffect) GetDescription() string {
	return "Shield no longer expires at the end of turn"
}