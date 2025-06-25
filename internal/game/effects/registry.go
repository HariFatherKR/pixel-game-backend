package effects

import (
	"fmt"
	"strings"
)

// EffectFactory is a function that creates a CardEffect instance
type EffectFactory func(params map[string]interface{}) (CardEffect, error)

// EffectRegistry manages all available card effects
type EffectRegistry struct {
	effects map[string]EffectFactory
}

// NewEffectRegistry creates a new effect registry
func NewEffectRegistry() *EffectRegistry {
	registry := &EffectRegistry{
		effects: make(map[string]EffectFactory),
	}
	
	// Register all effects
	registry.registerEffects()
	
	return registry
}

// registerEffects registers all available effects
func (r *EffectRegistry) registerEffects() {
	// Damage effects
	r.effects["damage"] = func(params map[string]interface{}) (CardEffect, error) {
		damage, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("damage value required")
		}
		return NewDamageEffect(int(damage)), nil
	}
	
	r.effects["multi_hit_damage"] = func(params map[string]interface{}) (CardEffect, error) {
		damage, ok := params["damage_per_hit"].(float64)
		if !ok {
			return nil, fmt.Errorf("damage_per_hit required")
		}
		hits, ok := params["hit_count"].(float64)
		if !ok {
			return nil, fmt.Errorf("hit_count required")
		}
		return NewMultiHitDamageEffect(int(damage), int(hits)), nil
	}
	
	r.effects["area_damage"] = func(params map[string]interface{}) (CardEffect, error) {
		damage, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("damage value required")
		}
		return NewAreaDamageEffect(int(damage)), nil
	}
	
	// Shield effects
	r.effects["shield"] = func(params map[string]interface{}) (CardEffect, error) {
		shield, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("shield value required")
		}
		return NewShieldEffect(int(shield)), nil
	}
	
	r.effects["reflect_shield"] = func(params map[string]interface{}) (CardEffect, error) {
		shield, ok := params["shield"].(float64)
		if !ok {
			return nil, fmt.Errorf("shield value required")
		}
		reflect, ok := params["reflect_percent"].(float64)
		if !ok {
			return nil, fmt.Errorf("reflect_percent required")
		}
		return NewReflectShieldEffect(int(shield), reflect), nil
	}
	
	r.effects["barricade"] = func(params map[string]interface{}) (CardEffect, error) {
		return NewBarricadeEffect(), nil
	}
	
	// Draw effects
	r.effects["draw"] = func(params map[string]interface{}) (CardEffect, error) {
		count, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("draw count required")
		}
		return NewDrawEffect(int(count)), nil
	}
	
	r.effects["scry"] = func(params map[string]interface{}) (CardEffect, error) {
		count, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("scry count required")
		}
		return NewScryEffect(int(count)), nil
	}
	
	r.effects["draw_to_hand_size"] = func(params map[string]interface{}) (CardEffect, error) {
		size, ok := params["target_size"].(float64)
		if !ok {
			return nil, fmt.Errorf("target_size required")
		}
		return NewDrawToHandSizeEffect(int(size)), nil
	}
	
	// Buff effects
	r.effects["strength"] = func(params map[string]interface{}) (CardEffect, error) {
		amount, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("strength amount required")
		}
		return NewStrengthEffect(int(amount)), nil
	}
	
	r.effects["dexterity"] = func(params map[string]interface{}) (CardEffect, error) {
		amount, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("dexterity amount required")
		}
		return NewDexterityEffect(int(amount)), nil
	}
	
	r.effects["vulnerable"] = func(params map[string]interface{}) (CardEffect, error) {
		duration, ok := params["duration"].(float64)
		if !ok {
			return nil, fmt.Errorf("vulnerable duration required")
		}
		return NewVulnerableEffect(int(duration)), nil
	}
	
	r.effects["weak"] = func(params map[string]interface{}) (CardEffect, error) {
		duration, ok := params["duration"].(float64)
		if !ok {
			return nil, fmt.Errorf("weak duration required")
		}
		target, ok := params["target"].(string)
		if !ok {
			target = "enemy"
		}
		return NewWeakEffect(int(duration), target), nil
	}
	
	r.effects["frail"] = func(params map[string]interface{}) (CardEffect, error) {
		duration, ok := params["duration"].(float64)
		if !ok {
			return nil, fmt.Errorf("frail duration required")
		}
		return NewFrailEffect(int(duration)), nil
	}
	
	// Special effects
	r.effects["energy_gain"] = func(params map[string]interface{}) (CardEffect, error) {
		amount, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("energy amount required")
		}
		return NewEnergyGainEffect(int(amount)), nil
	}
	
	r.effects["heal"] = func(params map[string]interface{}) (CardEffect, error) {
		amount, ok := params["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("heal amount required")
		}
		return NewHealEffect(int(amount)), nil
	}
	
	r.effects["exhaust"] = func(params map[string]interface{}) (CardEffect, error) {
		count, ok := params["count"].(float64)
		if !ok {
			count = 0
		}
		targetSelf, ok := params["target_self"].(bool)
		if !ok {
			targetSelf = false
		}
		return NewExhaustEffect(int(count), targetSelf), nil
	}
	
	r.effects["retain"] = func(params map[string]interface{}) (CardEffect, error) {
		cardID, ok := params["card_id"].(string)
		if !ok {
			return nil, fmt.Errorf("card_id required for retain")
		}
		return NewRetainEffect(cardID), nil
	}
	
	r.effects["double_play"] = func(params map[string]interface{}) (CardEffect, error) {
		return NewDoublePlayEffect(), nil
	}
}

// CreateEffect creates an effect instance from type and parameters
func (r *EffectRegistry) CreateEffect(effectType string, params map[string]interface{}) (CardEffect, error) {
	factory, exists := r.effects[strings.ToLower(effectType)]
	if !exists {
		return nil, fmt.Errorf("unknown effect type: %s", effectType)
	}
	
	return factory(params)
}

// GetAvailableEffects returns a list of all registered effect types
func (r *EffectRegistry) GetAvailableEffects() []string {
	effects := make([]string, 0, len(r.effects))
	for effectType := range r.effects {
		effects = append(effects, effectType)
	}
	return effects
}