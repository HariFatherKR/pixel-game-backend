package effects

import (
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// StrengthEffect increases damage dealt
type StrengthEffect struct {
	amount int
}

// NewStrengthEffect creates a strength effect
func NewStrengthEffect(amount int) *StrengthEffect {
	return &StrengthEffect{amount: amount}
}

// Execute applies strength
func (e *StrengthEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Add or increase strength power
	if power, exists := ctx.PlayerState.ActivePowers["strength"]; exists {
		power.Stacks += e.amount
		ctx.PlayerState.ActivePowers["strength"] = power
	} else {
		ctx.PlayerState.ActivePowers["strength"] = domain.PowerState{
			PowerID:     "strength",
			Name:        "Strength",
			Description: fmt.Sprintf("Increases attack damage by %d", e.amount),
			Stacks:      e.amount,
			Duration:    -1, // Permanent
		}
	}

	result.Messages = append(result.Messages, 
		fmt.Sprintf("Gained %d strength", e.amount))

	return result, nil
}

// CanExecute checks if strength can be applied
func (e *StrengthEffect) CanExecute(ctx *EffectContext) (bool, string) {
	return true, ""
}

// GetType returns the effect type
func (e *StrengthEffect) GetType() string {
	return "strength"
}

// GetDescription returns the effect description
func (e *StrengthEffect) GetDescription() string {
	return fmt.Sprintf("Gain %d strength", e.amount)
}

// DexterityEffect increases shield gained
type DexterityEffect struct {
	amount int
}

// NewDexterityEffect creates a dexterity effect
func NewDexterityEffect(amount int) *DexterityEffect {
	return &DexterityEffect{amount: amount}
}

// Execute applies dexterity
func (e *DexterityEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Add or increase dexterity power
	if power, exists := ctx.PlayerState.ActivePowers["dexterity"]; exists {
		power.Stacks += e.amount
		ctx.PlayerState.ActivePowers["dexterity"] = power
	} else {
		ctx.PlayerState.ActivePowers["dexterity"] = domain.PowerState{
			PowerID:     "dexterity",
			Name:        "Dexterity",
			Description: fmt.Sprintf("Increases shield gained by %d", e.amount),
			Stacks:      e.amount,
			Duration:    -1, // Permanent
		}
	}

	result.Messages = append(result.Messages, 
		fmt.Sprintf("Gained %d dexterity", e.amount))

	return result, nil
}

// CanExecute checks if dexterity can be applied
func (e *DexterityEffect) CanExecute(ctx *EffectContext) (bool, string) {
	return true, ""
}

// GetType returns the effect type
func (e *DexterityEffect) GetType() string {
	return "dexterity"
}

// GetDescription returns the effect description
func (e *DexterityEffect) GetDescription() string {
	return fmt.Sprintf("Gain %d dexterity", e.amount)
}

// VulnerableEffect makes target take more damage
type VulnerableEffect struct {
	duration int
}

// NewVulnerableEffect creates a vulnerable effect
func NewVulnerableEffect(duration int) *VulnerableEffect {
	return &VulnerableEffect{duration: duration}
}

// Execute applies vulnerable debuff
func (e *VulnerableEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	if ctx.TargetID == "" || ctx.EnemyState == nil {
		return result, fmt.Errorf("no valid target")
	}

	// Apply vulnerable to enemy
	debuff := domain.DebuffState{
		DebuffID:    "vulnerable",
		Name:        "Vulnerable",
		Description: "Takes 50% more damage",
		Value:       50,
		Duration:    e.duration,
	}

	// Check if already vulnerable
	found := false
	for i, existing := range ctx.EnemyState.Debuffs {
		if existing.DebuffID == "vulnerable" {
			ctx.EnemyState.Debuffs[i].Duration = e.duration
			found = true
			break
		}
	}

	if !found {
		ctx.EnemyState.Debuffs = append(ctx.EnemyState.Debuffs, debuff)
		result.DebuffsApplied = append(result.DebuffsApplied, debuff)
	}

	result.Messages = append(result.Messages, 
		fmt.Sprintf("Applied vulnerable for %d turns", e.duration))

	return result, nil
}

// CanExecute checks if vulnerable can be applied
func (e *VulnerableEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if ctx.TargetID == "" {
		return false, "no target selected"
	}
	if ctx.EnemyState == nil {
		return false, "invalid target"
	}
	return true, ""
}

// GetType returns the effect type
func (e *VulnerableEffect) GetType() string {
	return "vulnerable"
}

// GetDescription returns the effect description
func (e *VulnerableEffect) GetDescription() string {
	return fmt.Sprintf("Apply vulnerable for %d turns", e.duration)
}

// WeakEffect reduces damage dealt
type WeakEffect struct {
	duration int
	target   string // "player" or "enemy"
}

// NewWeakEffect creates a weak effect
func NewWeakEffect(duration int, target string) *WeakEffect {
	return &WeakEffect{
		duration: duration,
		target:   target,
	}
}

// Execute applies weak debuff
func (e *WeakEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	debuff := domain.DebuffState{
		DebuffID:    "weak",
		Name:        "Weak",
		Description: "Deals 25% less damage",
		Value:       25,
		Duration:    e.duration,
	}

	if e.target == "enemy" && ctx.EnemyState != nil {
		// Apply to enemy
		found := false
		for i, existing := range ctx.EnemyState.Debuffs {
			if existing.DebuffID == "weak" {
				ctx.EnemyState.Debuffs[i].Duration = e.duration
				found = true
				break
			}
		}
		if !found {
			ctx.EnemyState.Debuffs = append(ctx.EnemyState.Debuffs, debuff)
			result.DebuffsApplied = append(result.DebuffsApplied, debuff)
		}
		result.Messages = append(result.Messages, 
			fmt.Sprintf("Applied weak to enemy for %d turns", e.duration))
	} else if e.target == "player" {
		// Apply to player
		found := false
		for i, existing := range ctx.PlayerState.Debuffs {
			if existing.DebuffID == "weak" {
				ctx.PlayerState.Debuffs[i].Duration = e.duration
				found = true
				break
			}
		}
		if !found {
			ctx.PlayerState.Debuffs = append(ctx.PlayerState.Debuffs, debuff)
			result.DebuffsApplied = append(result.DebuffsApplied, debuff)
		}
		result.Messages = append(result.Messages, 
			fmt.Sprintf("Applied weak to player for %d turns", e.duration))
	}

	return result, nil
}

// CanExecute checks if weak can be applied
func (e *WeakEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if e.target == "enemy" && ctx.EnemyState == nil {
		return false, "no enemy target"
	}
	return true, ""
}

// GetType returns the effect type
func (e *WeakEffect) GetType() string {
	return "weak"
}

// GetDescription returns the effect description
func (e *WeakEffect) GetDescription() string {
	return fmt.Sprintf("Apply weak to %s for %d turns", e.target, e.duration)
}

// FrailEffect reduces shield gained
type FrailEffect struct {
	duration int
}

// NewFrailEffect creates a frail effect
func NewFrailEffect(duration int) *FrailEffect {
	return &FrailEffect{duration: duration}
}

// Execute applies frail debuff
func (e *FrailEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	debuff := domain.DebuffState{
		DebuffID:    "frail",
		Name:        "Frail",
		Description: "Gain 25% less shield",
		Value:       25,
		Duration:    e.duration,
	}

	// Apply to player
	found := false
	for i, existing := range ctx.PlayerState.Debuffs {
		if existing.DebuffID == "frail" {
			ctx.PlayerState.Debuffs[i].Duration = e.duration
			found = true
			break
		}
	}

	if !found {
		ctx.PlayerState.Debuffs = append(ctx.PlayerState.Debuffs, debuff)
		result.DebuffsApplied = append(result.DebuffsApplied, debuff)
	}

	result.Messages = append(result.Messages, 
		fmt.Sprintf("Applied frail for %d turns", e.duration))

	return result, nil
}

// CanExecute checks if frail can be applied
func (e *FrailEffect) CanExecute(ctx *EffectContext) (bool, string) {
	return true, ""
}

// GetType returns the effect type
func (e *FrailEffect) GetType() string {
	return "frail"
}

// GetDescription returns the effect description
func (e *FrailEffect) GetDescription() string {
	return fmt.Sprintf("Apply frail for %d turns", e.duration)
}