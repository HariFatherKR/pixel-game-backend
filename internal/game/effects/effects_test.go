package effects

import (
	"testing"
	"github.com/yourusername/pixel-game/internal/domain"
)

func TestDamageEffect(t *testing.T) {
	tests := []struct {
		name           string
		damage         int
		playerStrength int
		enemyVulnerable bool
		playerWeak     bool
		expectedDamage int
	}{
		{
			name:           "Basic damage",
			damage:         10,
			playerStrength: 0,
			enemyVulnerable: false,
			playerWeak:     false,
			expectedDamage: 10,
		},
		{
			name:           "Damage with strength",
			damage:         10,
			playerStrength: 5,
			enemyVulnerable: false,
			playerWeak:     false,
			expectedDamage: 15,
		},
		{
			name:           "Damage against vulnerable enemy",
			damage:         10,
			playerStrength: 0,
			enemyVulnerable: true,
			playerWeak:     false,
			expectedDamage: 15,
		},
		{
			name:           "Damage with weak player",
			damage:         10,
			playerStrength: 0,
			enemyVulnerable: false,
			playerWeak:     true,
			expectedDamage: 7,
		},
		{
			name:           "Combined modifiers",
			damage:         10,
			playerStrength: 5,
			enemyVulnerable: true,
			playerWeak:     true,
			expectedDamage: 16, // (10 + 5) * 1.5 * 0.75 = 16.875 -> 16
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			playerState := &domain.PlayerState{
				ActivePowers: make(map[string]domain.PowerState),
				Debuffs:      []domain.DebuffState{},
			}
			
			if tt.playerStrength > 0 {
				playerState.ActivePowers["strength"] = domain.PowerState{
					PowerID: "strength",
					Stacks:  tt.playerStrength,
				}
			}
			
			if tt.playerWeak {
				playerState.Debuffs = append(playerState.Debuffs, domain.DebuffState{
					DebuffID: "weak",
					Value:    25,
				})
			}

			enemyState := &domain.EnemyState{
				Name:      "Test Enemy",
				Health:    100,
				MaxHealth: 100,
				Shield:    0,
				Debuffs:   []domain.DebuffState{},
			}
			
			if tt.enemyVulnerable {
				enemyState.Debuffs = append(enemyState.Debuffs, domain.DebuffState{
					DebuffID: "vulnerable",
					Value:    50,
				})
			}

			ctx := &EffectContext{
				PlayerState: playerState,
				EnemyState:  enemyState,
				TargetID:    "enemy",
			}

			// Execute
			effect := NewDamageEffect(tt.damage)
			result, err := effect.Execute(ctx)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if result.Damage != tt.expectedDamage {
				t.Errorf("expected damage %d, got %d", tt.expectedDamage, result.Damage)
			}
			
			expectedHealth := 100 - tt.expectedDamage
			if enemyState.Health != expectedHealth {
				t.Errorf("expected enemy health %d, got %d", expectedHealth, enemyState.Health)
			}
		})
	}
}

func TestShieldEffect(t *testing.T) {
	tests := []struct {
		name           string
		shield         int
		playerDexterity int
		playerFrail    bool
		expectedShield int
	}{
		{
			name:           "Basic shield",
			shield:         10,
			playerDexterity: 0,
			playerFrail:    false,
			expectedShield: 10,
		},
		{
			name:           "Shield with dexterity",
			shield:         10,
			playerDexterity: 5,
			playerFrail:    false,
			expectedShield: 15,
		},
		{
			name:           "Shield with frail",
			shield:         10,
			playerDexterity: 0,
			playerFrail:    true,
			expectedShield: 7,
		},
		{
			name:           "Shield with dexterity and frail",
			shield:         10,
			playerDexterity: 5,
			playerFrail:    true,
			expectedShield: 11, // (10 + 5) * 0.75 = 11.25 -> 11
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			playerState := &domain.PlayerState{
				Shield:       0,
				ActivePowers: make(map[string]domain.PowerState),
				Debuffs:      []domain.DebuffState{},
			}
			
			if tt.playerDexterity > 0 {
				playerState.ActivePowers["dexterity"] = domain.PowerState{
					PowerID: "dexterity",
					Stacks:  tt.playerDexterity,
				}
			}
			
			if tt.playerFrail {
				playerState.Debuffs = append(playerState.Debuffs, domain.DebuffState{
					DebuffID: "frail",
					Value:    25,
				})
			}

			ctx := &EffectContext{
				PlayerState: playerState,
			}

			// Execute
			effect := NewShieldEffect(tt.shield)
			result, err := effect.Execute(ctx)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if result.ShieldGained != tt.expectedShield {
				t.Errorf("expected shield %d, got %d", tt.expectedShield, result.ShieldGained)
			}
			
			if playerState.Shield != tt.expectedShield {
				t.Errorf("expected player shield %d, got %d", tt.expectedShield, playerState.Shield)
			}
		})
	}
}

func TestDrawEffect(t *testing.T) {
	tests := []struct {
		name          string
		drawCount     int
		handSize      int
		drawPileSize  int
		discardSize   int
		expectedDrawn int
	}{
		{
			name:          "Basic draw",
			drawCount:     3,
			handSize:      5,
			drawPileSize:  10,
			discardSize:   0,
			expectedDrawn: 3,
		},
		{
			name:          "Draw with full hand",
			drawCount:     3,
			handSize:      10,
			drawPileSize:  10,
			discardSize:   0,
			expectedDrawn: 0,
		},
		{
			name:          "Draw with nearly full hand",
			drawCount:     3,
			handSize:      8,
			drawPileSize:  10,
			discardSize:   0,
			expectedDrawn: 2,
		},
		{
			name:          "Draw with empty draw pile",
			drawCount:     3,
			handSize:      5,
			drawPileSize:  0,
			discardSize:   5,
			expectedDrawn: 3,
		},
		{
			name:          "Draw more than available",
			drawCount:     10,
			handSize:      5,
			drawPileSize:  3,
			discardSize:   0,
			expectedDrawn: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			hand := make([]string, tt.handSize)
			for i := 0; i < tt.handSize; i++ {
				hand[i] = string(rune('a' + i))
			}
			
			drawPile := make([]string, tt.drawPileSize)
			for i := 0; i < tt.drawPileSize; i++ {
				drawPile[i] = string(rune('A' + i))
			}
			
			discardPile := make([]string, tt.discardSize)
			for i := 0; i < tt.discardSize; i++ {
				discardPile[i] = string(rune('1' + i))
			}

			playerState := &domain.PlayerState{
				Hand:        hand,
				DrawPile:    drawPile,
				DiscardPile: discardPile,
			}

			ctx := &EffectContext{
				PlayerState: playerState,
			}

			// Execute
			effect := NewDrawEffect(tt.drawCount)
			result, err := effect.Execute(ctx)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if len(result.CardsDrawn) != tt.expectedDrawn {
				t.Errorf("expected to draw %d cards, got %d", tt.expectedDrawn, len(result.CardsDrawn))
			}
			
			expectedHandSize := tt.handSize + tt.expectedDrawn
			if len(playerState.Hand) != expectedHandSize {
				t.Errorf("expected hand size %d, got %d", expectedHandSize, len(playerState.Hand))
			}
		})
	}
}

func TestEffectRegistry(t *testing.T) {
	registry := NewEffectRegistry()
	
	// Test creating damage effect
	params := map[string]interface{}{
		"value": float64(10),
	}
	
	effect, err := registry.CreateEffect("damage", params)
	if err != nil {
		t.Errorf("failed to create damage effect: %v", err)
	}
	
	if effect.GetType() != "damage" {
		t.Errorf("expected effect type 'damage', got '%s'", effect.GetType())
	}
	
	// Test creating unknown effect
	_, err = registry.CreateEffect("unknown_effect", params)
	if err == nil {
		t.Error("expected error for unknown effect type")
	}
}

func TestEffectExecutor(t *testing.T) {
	executor := NewExecutor()
	
	// Create a test card
	card := &domain.Card{
		ID:     "test_card",
		Name:   "Test Card",
		Type:   domain.CardTypeAction,
		Cost:   2,
		Effects: []byte(`[
			{"type": "damage", "target": "enemy", "value": 10},
			{"type": "shield", "target": "self", "value": 5}
		]`),
	}
	
	// Setup game state
	playerState := &domain.PlayerState{
		Health:       100,
		MaxHealth:    100,
		Shield:       0,
		ActivePowers: make(map[string]domain.PowerState),
		Debuffs:      []domain.DebuffState{},
	}
	
	enemyState := &domain.EnemyState{
		Name:      "Test Enemy",
		Health:    50,
		MaxHealth: 50,
		Shield:    0,
		Debuffs:   []domain.DebuffState{},
	}
	
	gameState := &domain.GameState{}
	targetID := "enemy"
	
	// Execute card effects
	result, err := executor.ExecuteCardEffects(card, playerState, enemyState, gameState, &targetID)
	if err != nil {
		t.Errorf("failed to execute card effects: %v", err)
	}
	
	// Verify results
	if result.DamageDealt != 10 {
		t.Errorf("expected damage dealt 10, got %d", result.DamageDealt)
	}
	
	if result.ShieldGained != 5 {
		t.Errorf("expected shield gained 5, got %d", result.ShieldGained)
	}
	
	if enemyState.Health != 40 {
		t.Errorf("expected enemy health 40, got %d", enemyState.Health)
	}
	
	if playerState.Shield != 5 {
		t.Errorf("expected player shield 5, got %d", playerState.Shield)
	}
}