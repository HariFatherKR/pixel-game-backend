package ai

import (
	"testing"
	"github.com/yourusername/pixel-game/internal/domain"
)

func TestAggressiveAI(t *testing.T) {
	ai := NewAggressiveAI(10, 1.0)
	
	// 기본 테스트 컨텍스트 생성
	ctx := &AIContext{
		EnemyState: &domain.EnemyState{
			ID:           "test_enemy",
			Name:         "테스트 적",
			Health:       50,
			MaxHealth:    50,
			Shield:       0,
			Buffs:        []domain.BuffState{},
			Debuffs:      []domain.DebuffState{},
			Intent: domain.EnemyIntent{
				Type:        "ATTACK",
				Value:       10,
				Description: "공격 준비 중",
			},
		},
		PlayerState: &domain.PlayerState{
			Health:       80,
			MaxHealth:    80,
			Shield:       0,
			ActivePowers: make(map[string]domain.PowerState),
			Buffs:        []domain.BuffState{},
			Debuffs:      []domain.DebuffState{},
		},
		GameState:   &domain.GameState{},
		TurnNumber:  1,
		FloorNumber: 1,
	}
	
	t.Run("AI 기본 정보 확인", func(t *testing.T) {
		if ai.GetName() != "Aggressive" {
			t.Errorf("AI 이름이 잘못됨: expected 'Aggressive', got '%s'", ai.GetName())
		}
		
		if ai.GetBehaviorType() != "AGGRESSIVE" {
			t.Errorf("AI 행동 유형이 잘못됨: expected 'AGGRESSIVE', got '%s'", ai.GetBehaviorType())
		}
	})
	
	t.Run("의도 계산 테스트", func(t *testing.T) {
		intent, err := ai.CalculateIntent(ctx)
		if err != nil {
			t.Errorf("의도 계산 중 오류: %v", err)
		}
		
		if intent == nil {
			t.Error("의도가 nil입니다")
		}
		
		// 공격적 AI는 주로 공격 의도를 가져야 함
		if intent.Type != "ATTACK" && intent.Type != "SPECIAL_ATTACK" && intent.Type != "BUFF" {
			t.Errorf("예상하지 못한 의도 타입: %s", intent.Type)
		}
	})
	
	t.Run("공격 실행 테스트", func(t *testing.T) {
		// 공격 의도로 설정
		ctx.EnemyState.Intent = domain.EnemyIntent{
			Type:        "ATTACK",
			Value:       12,
			Description: "12 데미지 공격",
		}
		
		initialHealth := ctx.PlayerState.Health
		result, err := ai.ExecuteAction(ctx)
		
		if err != nil {
			t.Errorf("공격 실행 중 오류: %v", err)
		}
		
		if !result.Success {
			t.Error("공격이 실패했습니다")
		}
		
		if result.Damage <= 0 {
			t.Error("데미지가 0 이하입니다")
		}
		
		if ctx.PlayerState.Health >= initialHealth {
			t.Error("플레이어 체력이 감소하지 않았습니다")
		}
	})
	
	t.Run("버프 상태에서 데미지 증가 테스트", func(t *testing.T) {
		// 분노 버프 추가
		ctx.EnemyState.Buffs = []domain.BuffState{
			{
				BuffID:   "rage",
				Name:     "분노",
				Value:    3,
				Duration: 2,
			},
		}
		
		intent, err := ai.CalculateIntent(ctx)
		if err != nil {
			t.Errorf("버프 상태에서 의도 계산 중 오류: %v", err)
		}
		
		// 버프가 있을 때 데미지가 증가해야 함
		if intent.Type == "ATTACK" && intent.Value <= 10 {
			t.Error("버프 상태에서 데미지가 증가하지 않았습니다")
		}
	})
}

func TestDefensiveAI(t *testing.T) {
	ai := NewDefensiveAI(8, 10, 8)
	
	ctx := &AIContext{
		EnemyState: &domain.EnemyState{
			ID:        "test_defensive",
			Name:      "방어적 적",
			Health:    20, // 낮은 체력으로 설정
			MaxHealth: 50,
			Shield:    0,
			Buffs:     []domain.BuffState{},
			Debuffs:   []domain.DebuffState{},
		},
		PlayerState: &domain.PlayerState{
			Health:       80,
			MaxHealth:    80,
			Shield:       0,
			ActivePowers: make(map[string]domain.PowerState),
			Buffs:        []domain.BuffState{},
			Debuffs:      []domain.DebuffState{},
		},
		GameState:   &domain.GameState{},
		TurnNumber:  1,
		FloorNumber: 1,
	}
	
	t.Run("낮은 체력에서 회복 의도 테스트", func(t *testing.T) {
		intent, err := ai.CalculateIntent(ctx)
		if err != nil {
			t.Errorf("의도 계산 중 오류: %v", err)
		}
		
		// 체력이 낮을 때는 회복이나 방어 의도를 가져야 함
		if intent.Type != "HEAL" && intent.Type != "DEFEND" {
			t.Errorf("낮은 체력에서 예상되는 의도가 아님: %s", intent.Type)
		}
	})
	
	t.Run("방어 실행 테스트", func(t *testing.T) {
		ctx.EnemyState.Intent = domain.EnemyIntent{
			Type:        "DEFEND",
			Value:       10,
			Description: "방어막 생성",
		}
		
		initialShield := ctx.EnemyState.Shield
		result, err := ai.ExecuteAction(ctx)
		
		if err != nil {
			t.Errorf("방어 실행 중 오류: %v", err)
		}
		
		if !result.Success {
			t.Error("방어가 실패했습니다")
		}
		
		if ctx.EnemyState.Shield <= initialShield {
			t.Error("방어막이 증가하지 않았습니다")
		}
	})
	
	t.Run("회복 실행 테스트", func(t *testing.T) {
		ctx.EnemyState.Intent = domain.EnemyIntent{
			Type:        "HEAL",
			Value:       8,
			Description: "체력 회복",
		}
		
		initialHealth := ctx.EnemyState.Health
		result, err := ai.ExecuteAction(ctx)
		
		if err != nil {
			t.Errorf("회복 실행 중 오류: %v", err)
		}
		
		if !result.Success {
			t.Error("회복이 실패했습니다")
		}
		
		if ctx.EnemyState.Health <= initialHealth {
			t.Error("체력이 회복되지 않았습니다")
		}
	})
}

func TestBalancedAI(t *testing.T) {
	ai := NewBalancedAI(10, 8, 6)
	
	ctx := &AIContext{
		EnemyState: &domain.EnemyState{
			ID:        "test_balanced",
			Name:      "균형 적",
			Health:    40,
			MaxHealth: 50,
			Shield:    0,
			Buffs:     []domain.BuffState{},
			Debuffs:   []domain.DebuffState{},
		},
		PlayerState: &domain.PlayerState{
			Health:       40, // 중간 체력
			MaxHealth:    80,
			Shield:       0,
			ActivePowers: make(map[string]domain.PowerState),
			Buffs:        []domain.BuffState{},
			Debuffs:      []domain.DebuffState{},
		},
		GameState:   &domain.GameState{},
		TurnNumber:  1,
		FloorNumber: 1,
	}
	
	t.Run("상황 분석 테스트", func(t *testing.T) {
		intent, err := ai.CalculateIntent(ctx)
		if err != nil {
			t.Errorf("의도 계산 중 오류: %v", err)
		}
		
		// 균형 AI는 모든 타입의 의도를 가질 수 있음
		validTypes := []string{"ATTACK", "DEFEND", "HEAL", "BUFF", "DEBUFF"}
		found := false
		for _, validType := range validTypes {
			if intent.Type == validType {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("예상되지 않은 의도 타입: %s", intent.Type)
		}
	})
	
	t.Run("디버프 실행 테스트", func(t *testing.T) {
		ctx.EnemyState.Intent = domain.EnemyIntent{
			Type:        "DEBUFF",
			Value:       50,
			Description: "플레이어 약화",
		}
		
		initialDebuffCount := len(ctx.PlayerState.Debuffs)
		result, err := ai.ExecuteAction(ctx)
		
		if err != nil {
			t.Errorf("디버프 실행 중 오류: %v", err)
		}
		
		if !result.Success {
			t.Error("디버프가 실패했습니다")
		}
		
		if len(ctx.PlayerState.Debuffs) <= initialDebuffCount {
			t.Error("플레이어에게 디버프가 적용되지 않았습니다")
		}
	})
}

func TestAIManager(t *testing.T) {
	manager := NewAIManager()
	
	t.Run("AI 등록 및 조회 테스트", func(t *testing.T) {
		// 기본 AI들이 등록되어 있는지 확인
		ais := []string{"aggressive", "defensive", "balanced"}
		
		for _, aiName := range ais {
			ai, err := manager.GetAI(aiName)
			if err != nil {
				t.Errorf("AI '%s' 조회 실패: %v", aiName, err)
			}
			
			if ai == nil {
				t.Errorf("AI '%s'가 nil입니다", aiName)
			}
		}
	})
	
	t.Run("적 타입별 AI 선택 테스트", func(t *testing.T) {
		testCases := []struct {
			enemyType string
			expected  string
		}{
			{"BASIC_ENEMY", "balanced"},
			{"BRUTE", "aggressive"},
			{"GUARDIAN", "defensive"},
			{"ELITE", "balanced"},
		}
		
		for _, tc := range testCases {
			ai, err := manager.SelectAIForEnemy(tc.enemyType, 1)
			if err != nil {
				t.Errorf("적 타입 '%s'에 대한 AI 선택 실패: %v", tc.enemyType, err)
			}
			
			expectedAI, _ := manager.GetAI(tc.expected)
			if ai.GetBehaviorType() != expectedAI.GetBehaviorType() {
				t.Errorf("잘못된 AI 선택: expected %s, got %s", 
					expectedAI.GetBehaviorType(), ai.GetBehaviorType())
			}
		}
	})
	
	t.Run("적 턴 처리 테스트", func(t *testing.T) {
		enemyState := &domain.EnemyState{
			ID:        "test_enemy",
			Name:      "테스트 적",
			Health:    50,
			MaxHealth: 50,
			Shield:    0,
			Buffs:     []domain.BuffState{},
			Debuffs:   []domain.DebuffState{},
			Intent: domain.EnemyIntent{
				Type:        "ATTACK",
				Value:       10,
				Description: "공격 준비 중",
			},
		}
		
		playerState := &domain.PlayerState{
			Health:       80,
			MaxHealth:    80,
			Shield:       0,
			ActivePowers: make(map[string]domain.PowerState),
			Buffs:        []domain.BuffState{},
			Debuffs:      []domain.DebuffState{},
		}
		
		gameState := &domain.GameState{}
		
		result, err := manager.ProcessEnemyTurn(
			enemyState,
			playerState,
			gameState,
			1, // turnNumber
			1, // floorNumber
			"aggressive",
		)
		
		if err != nil {
			t.Errorf("적 턴 처리 중 오류: %v", err)
		}
		
		if result == nil {
			t.Error("결과가 nil입니다")
		}
		
		if !result.Success {
			t.Error("적 턴 처리가 실패했습니다")
		}
	})
}

func TestAIValidation(t *testing.T) {
	manager := NewAIManager()
	
	ctx := &AIContext{
		EnemyState: &domain.EnemyState{
			ID:        "test_enemy",
			Health:    50,
			MaxHealth: 50,
			Shield:    0,
			Buffs:     []domain.BuffState{},
			Debuffs:   []domain.DebuffState{},
		},
		PlayerState: &domain.PlayerState{
			Health:       80,
			MaxHealth:    80,
			Shield:       0,
			ActivePowers: make(map[string]domain.PowerState),
			Buffs:        []domain.BuffState{},
			Debuffs:      []domain.DebuffState{},
		},
		GameState:   &domain.GameState{},
		TurnNumber:  1,
		FloorNumber: 1,
	}
	
	t.Run("행동 유효성 검사 테스트", func(t *testing.T) {
		canExecute, reason, err := manager.ValidateAIAction(
			ctx.EnemyState,
			ctx.PlayerState,
			ctx.GameState,
			1, 1,
			"aggressive",
			"ATTACK",
		)
		
		if err != nil {
			t.Errorf("행동 유효성 검사 중 오류: %v", err)
		}
		
		if !canExecute {
			t.Errorf("공격 행동이 유효하지 않음: %s", reason)
		}
	})
	
	t.Run("최대 체력에서 회복 불가 테스트", func(t *testing.T) {
		ctx.EnemyState.Health = ctx.EnemyState.MaxHealth
		
		canExecute, reason, err := manager.ValidateAIAction(
			ctx.EnemyState,
			ctx.PlayerState,
			ctx.GameState,
			1, 1,
			"defensive",
			"HEAL",
		)
		
		if err != nil {
			t.Errorf("회복 유효성 검사 중 오류: %v", err)
		}
		
		if canExecute {
			t.Error("최대 체력에서 회복이 가능하다고 판단됨")
		}
		
		if reason == "" {
			t.Error("유효하지 않은 이유가 제공되지 않음")
		}
	})
}