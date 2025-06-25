package effects

import (
	"encoding/json"
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// Executor handles the execution of card effects
type Executor struct {
	registry *EffectRegistry
}

// NewExecutor creates a new effect executor
func NewExecutor() *Executor {
	return &Executor{
		registry: NewEffectRegistry(),
	}
}

// ExecuteCardEffects executes all effects from a card
func (e *Executor) ExecuteCardEffects(
	card *domain.Card,
	playerState *domain.PlayerState,
	enemyState *domain.EnemyState,
	gameState *domain.GameState,
	targetID *string,
) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Success:        true,
		DamageDealt:    0,
		HealingDone:    0,
		ShieldGained:   0,
		CardsDrawn:     []string{},
		BuffsApplied:   []domain.BuffState{},
		DebuffsApplied: []domain.DebuffState{},
		Messages:       []string{},
	}

	// Create effect context
	ctx := &EffectContext{
		PlayerState: playerState,
		EnemyState:  enemyState,
		GameState:   gameState,
		SourceCard:  card,
		TargetID:    "",
	}

	if targetID != nil {
		ctx.TargetID = *targetID
	}

	// Get card effects
	cardEffects, err := card.GetEffects()
	if err != nil {
		return result, fmt.Errorf("failed to parse card effects: %w", err)
	}

	// Parse and execute each effect
	for _, effectData := range cardEffects {
		effectResult, err := e.executeEffect(effectData, ctx)
		if err != nil {
			return result, fmt.Errorf("failed to execute effect: %w", err)
		}

		// Merge results
		result.mergeWith(effectResult)
	}

	return result, nil
}

// executeEffect executes a single effect
func (e *Executor) executeEffect(effectData domain.CardEffect, ctx *EffectContext) (*EffectResult, error) {
	// Parse effect parameters
	params := make(map[string]interface{})
	params["value"] = float64(effectData.Value)
	params["target"] = effectData.Target

	// Parse additional parameters from JSON if present
	if len(effectData.Parameters) > 0 {
		var additionalParams map[string]interface{}
		if err := json.Unmarshal(effectData.Parameters, &additionalParams); err == nil {
			for k, v := range additionalParams {
				params[k] = v
			}
		}
	}

	// Create effect instance
	effect, err := e.registry.CreateEffect(effectData.Type, params)
	if err != nil {
		return nil, err
	}

	// Check if effect can be executed
	if canExecute, reason := effect.CanExecute(ctx); !canExecute {
		return &EffectResult{
			Success:  false,
			Messages: []string{fmt.Sprintf("Cannot execute %s: %s", effect.GetType(), reason)},
		}, nil
	}

	// Execute the effect
	return effect.Execute(ctx)
}

// ExecutionResult represents the aggregated result of card execution
type ExecutionResult struct {
	Success        bool                  `json:"success"`
	DamageDealt    int                   `json:"damage_dealt,omitempty"`
	HealingDone    int                   `json:"healing_done,omitempty"`
	ShieldGained   int                   `json:"shield_gained,omitempty"`
	CardsDrawn     []string              `json:"cards_drawn,omitempty"`
	BuffsApplied   []domain.BuffState    `json:"buffs_applied,omitempty"`
	DebuffsApplied []domain.DebuffState  `json:"debuffs_applied,omitempty"`
	Messages       []string              `json:"messages"`
}

// mergeWith merges another effect result into this one
func (r *ExecutionResult) mergeWith(other *EffectResult) {
	r.DamageDealt += other.Damage
	r.HealingDone += other.Healing
	r.ShieldGained += other.ShieldGained
	r.CardsDrawn = append(r.CardsDrawn, other.CardsDrawn...)
	r.BuffsApplied = append(r.BuffsApplied, other.BuffsApplied...)
	r.DebuffsApplied = append(r.DebuffsApplied, other.DebuffsApplied...)
	r.Messages = append(r.Messages, other.Messages...)
	if !other.Success {
		r.Success = false
	}
}

// ToMap converts the result to a map for JSON response
func (r *ExecutionResult) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"success":  r.Success,
		"messages": r.Messages,
	}

	if r.DamageDealt > 0 {
		result["damage_dealt"] = r.DamageDealt
	}
	if r.HealingDone > 0 {
		result["healing_done"] = r.HealingDone
	}
	if r.ShieldGained > 0 {
		result["shield_gained"] = r.ShieldGained
	}
	if len(r.CardsDrawn) > 0 {
		result["cards_drawn"] = r.CardsDrawn
	}
	if len(r.BuffsApplied) > 0 {
		result["buffs_applied"] = r.BuffsApplied
	}
	if len(r.DebuffsApplied) > 0 {
		result["debuffs_applied"] = r.DebuffsApplied
	}

	return result
}