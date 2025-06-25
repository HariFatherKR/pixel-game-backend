package effects

import (
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// DamageEffect implements damage dealing
type DamageEffect struct {
	baseDamage int
}

// NewDamageEffect creates a new damage effect
func NewDamageEffect(damage int) *DamageEffect {
	return &DamageEffect{
		baseDamage: damage,
	}
}

// Execute deals damage to the target
func (e *DamageEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Calculate actual damage (can be modified by buffs/debuffs)
	damage := e.calculateDamage(ctx)

	// Apply damage to enemy
	if ctx.TargetID != "" && ctx.EnemyState != nil {
		actualDamage := e.applyDamageToEnemy(ctx.EnemyState, damage)
		result.Damage = actualDamage
		result.Messages = append(result.Messages, 
			fmt.Sprintf("Dealt %d damage to %s", actualDamage, ctx.EnemyState.Name))
	}

	return result, nil
}

// CanExecute checks if damage can be dealt
func (e *DamageEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if ctx.TargetID == "" {
		return false, "no target selected"
	}
	if ctx.EnemyState == nil {
		return false, "invalid target"
	}
	return true, ""
}

// GetType returns the effect type
func (e *DamageEffect) GetType() string {
	return "damage"
}

// GetDescription returns the effect description
func (e *DamageEffect) GetDescription() string {
	return fmt.Sprintf("Deal %d damage", e.baseDamage)
}

// calculateDamage calculates actual damage considering modifiers
func (e *DamageEffect) calculateDamage(ctx *EffectContext) int {
	damage := e.baseDamage

	// Check for strength buff on player
	if strength, exists := ctx.PlayerState.ActivePowers["strength"]; exists {
		damage += strength.Stacks
	}

	// Check for vulnerable debuff on enemy
	for _, debuff := range ctx.EnemyState.Debuffs {
		if debuff.DebuffID == "vulnerable" {
			damage = int(float64(damage) * 1.5)
			break
		}
	}

	// Check for weak debuff on player
	for _, debuff := range ctx.PlayerState.Debuffs {
		if debuff.DebuffID == "weak" {
			damage = int(float64(damage) * 0.75)
			break
		}
	}

	return damage
}

// applyDamageToEnemy applies damage to enemy considering shields
func (e *DamageEffect) applyDamageToEnemy(enemy *domain.EnemyState, damage int) int {
	actualDamage := damage

	// Apply to shield first
	if enemy.Shield > 0 {
		if enemy.Shield >= damage {
			enemy.Shield -= damage
			return damage
		}
		actualDamage = damage - enemy.Shield
		enemy.Shield = 0
	}

	// Apply remaining damage to health
	enemy.Health -= actualDamage
	if enemy.Health < 0 {
		enemy.Health = 0
	}

	return damage
}

// MultiHitDamageEffect implements multiple hits
type MultiHitDamageEffect struct {
	damagePerHit int
	hitCount     int
}

// NewMultiHitDamageEffect creates a multi-hit damage effect
func NewMultiHitDamageEffect(damagePerHit, hitCount int) *MultiHitDamageEffect {
	return &MultiHitDamageEffect{
		damagePerHit: damagePerHit,
		hitCount:     hitCount,
	}
}

// Execute deals damage multiple times
func (e *MultiHitDamageEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	totalDamage := 0
	baseEffect := NewDamageEffect(e.damagePerHit)

	for i := 0; i < e.hitCount; i++ {
		hitResult, err := baseEffect.Execute(ctx)
		if err != nil {
			return result, err
		}
		totalDamage += hitResult.Damage
	}

	result.Damage = totalDamage
	result.Messages = append(result.Messages, 
		fmt.Sprintf("Dealt %d damage over %d hits", totalDamage, e.hitCount))

	return result, nil
}

// CanExecute checks if damage can be dealt
func (e *MultiHitDamageEffect) CanExecute(ctx *EffectContext) (bool, string) {
	baseEffect := NewDamageEffect(e.damagePerHit)
	return baseEffect.CanExecute(ctx)
}

// GetType returns the effect type
func (e *MultiHitDamageEffect) GetType() string {
	return "multi_hit_damage"
}

// GetDescription returns the effect description
func (e *MultiHitDamageEffect) GetDescription() string {
	return fmt.Sprintf("Deal %d damage %d times", e.damagePerHit, e.hitCount)
}

// AreaDamageEffect implements damage to all enemies
type AreaDamageEffect struct {
	damage int
}

// NewAreaDamageEffect creates an area damage effect
func NewAreaDamageEffect(damage int) *AreaDamageEffect {
	return &AreaDamageEffect{
		damage: damage,
	}
}

// Execute deals damage to all enemies
func (e *AreaDamageEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// In the current implementation, we only have one enemy
	// In future, this would iterate over all enemies
	if ctx.EnemyState != nil {
		damageEffect := NewDamageEffect(e.damage)
		damageResult, err := damageEffect.Execute(ctx)
		if err != nil {
			return result, err
		}
		result.Damage = damageResult.Damage
		result.Messages = append(result.Messages, 
			fmt.Sprintf("Dealt %d damage to all enemies", damageResult.Damage))
	}

	return result, nil
}

// CanExecute checks if area damage can be dealt
func (e *AreaDamageEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if ctx.EnemyState == nil {
		return false, "no enemies present"
	}
	return true, ""
}

// GetType returns the effect type
func (e *AreaDamageEffect) GetType() string {
	return "area_damage"
}

// GetDescription returns the effect description
func (e *AreaDamageEffect) GetDescription() string {
	return fmt.Sprintf("Deal %d damage to all enemies", e.damage)
}