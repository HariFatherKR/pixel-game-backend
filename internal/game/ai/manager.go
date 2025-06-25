package ai

import (
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// AIManager AI 시스템을 관리하는 매니저
type AIManager struct {
	registry *AIRegistry
}

// NewAIManager 새로운 AI 매니저 생성
func NewAIManager() *AIManager {
	manager := &AIManager{
		registry: NewAIRegistry(),
	}
	
	// 기본 AI들 등록
	manager.registerDefaultAIs()
	
	return manager
}

// registerDefaultAIs 기본 AI들을 등록
func (m *AIManager) registerDefaultAIs() {
	// 공격적 AI
	aggressiveAI := NewAggressiveAI(12, 1.5) // 기본 데미지 12, 층당 1.5씩 증가
	m.registry.Register("aggressive", aggressiveAI)
	
	// 방어적 AI  
	defensiveAI := NewDefensiveAI(8, 10, 8) // 데미지 8, 방어막 10, 회복 8
	m.registry.Register("defensive", defensiveAI)
	
	// 균형 AI
	balancedAI := NewBalancedAI(10, 8, 6) // 데미지 10, 방어막 8, 회복 6
	m.registry.Register("balanced", balancedAI)
}

// GetAI AI 이름으로 AI 인스턴스 가져오기
func (m *AIManager) GetAI(aiName string) (EnemyAI, error) {
	ai, exists := m.registry.Get(aiName)
	if !exists {
		return nil, fmt.Errorf("AI '%s'를 찾을 수 없습니다", aiName)
	}
	return ai, nil
}

// SelectAIForEnemy 적의 타입에 따라 적절한 AI 선택
func (m *AIManager) SelectAIForEnemy(enemyType string, floorNumber int) (EnemyAI, error) {
	switch enemyType {
	case "BASIC_ENEMY":
		// 기본 적은 균형 AI 사용
		return m.GetAI("balanced")
	case "BRUTE":
		// 무력형 적은 공격적 AI 사용
		return m.GetAI("aggressive")
	case "GUARDIAN":
		// 수호자형 적은 방어적 AI 사용
		return m.GetAI("defensive")
	case "ELITE":
		// 엘리트 적은 층수에 따라 선택
		if floorNumber <= 3 {
			return m.GetAI("balanced")
		} else if floorNumber <= 6 {
			return m.GetAI("aggressive")
		} else {
			return m.GetAI("defensive")
		}
	default:
		// 기본값은 균형 AI
		return m.GetAI("balanced")
	}
}

// ProcessEnemyTurn 적의 턴을 처리
func (m *AIManager) ProcessEnemyTurn(
	enemyState *domain.EnemyState,
	playerState *domain.PlayerState,
	gameState *domain.GameState,
	turnNumber int,
	floorNumber int,
	aiName string,
) (*AIResult, error) {
	// AI 가져오기
	ai, err := m.GetAI(aiName)
	if err != nil {
		return nil, fmt.Errorf("AI 처리 중 오류: %w", err)
	}
	
	// AI 컨텍스트 생성
	ctx := &AIContext{
		EnemyState:  enemyState,
		PlayerState: playerState,
		GameState:   gameState,
		TurnNumber:  turnNumber,
		FloorNumber: floorNumber,
	}
	
	// 현재 의도에 따라 행동 실행
	result, err := ai.ExecuteAction(ctx)
	if err != nil {
		return nil, fmt.Errorf("AI 행동 실행 중 오류: %w", err)
	}
	
	// 버프/디버프 지속시간 감소 처리
	m.updateBuffsAndDebuffs(enemyState, playerState)
	
	return result, nil
}

// CalculateNextIntent 다음 턴 의도 계산
func (m *AIManager) CalculateNextIntent(
	enemyState *domain.EnemyState,
	playerState *domain.PlayerState,
	gameState *domain.GameState,
	turnNumber int,
	floorNumber int,
	aiName string,
) (*domain.EnemyIntent, error) {
	// AI 가져오기
	ai, err := m.GetAI(aiName)
	if err != nil {
		return nil, fmt.Errorf("AI 의도 계산 중 오류: %w", err)
	}
	
	// AI 컨텍스트 생성
	ctx := &AIContext{
		EnemyState:  enemyState,
		PlayerState: playerState,
		GameState:   gameState,
		TurnNumber:  turnNumber + 1, // 다음 턴
		FloorNumber: floorNumber,
	}
	
	// 다음 의도 계산
	intent, err := ai.CalculateIntent(ctx)
	if err != nil {
		return nil, fmt.Errorf("AI 의도 계산 중 오류: %w", err)
	}
	
	return intent, nil
}

// updateBuffsAndDebuffs 버프와 디버프의 지속시간을 업데이트
func (m *AIManager) updateBuffsAndDebuffs(enemyState *domain.EnemyState, playerState *domain.PlayerState) {
	// 적의 버프 업데이트
	enemyState.Buffs = m.updateBuffsList(enemyState.Buffs)
	
	// 적의 디버프 업데이트
	enemyState.Debuffs = m.updateDebuffsList(enemyState.Debuffs)
	
	// 플레이어의 디버프 업데이트 (적이 적용한 것들)
	playerState.Debuffs = m.updateDebuffsList(playerState.Debuffs)
}

// updateBuffsList 버프 목록 업데이트
func (m *AIManager) updateBuffsList(buffs []domain.BuffState) []domain.BuffState {
	var updatedBuffs []domain.BuffState
	
	for _, buff := range buffs {
		if buff.Duration > 0 {
			buff.Duration--
			if buff.Duration > 0 {
				updatedBuffs = append(updatedBuffs, buff)
			}
		} else if buff.Duration == -1 { // 영구 버프
			updatedBuffs = append(updatedBuffs, buff)
		}
	}
	
	return updatedBuffs
}

// updateDebuffsList 디버프 목록 업데이트
func (m *AIManager) updateDebuffsList(debuffs []domain.DebuffState) []domain.DebuffState {
	var updatedDebuffs []domain.DebuffState
	
	for _, debuff := range debuffs {
		if debuff.Duration > 0 {
			debuff.Duration--
			if debuff.Duration > 0 {
				updatedDebuffs = append(updatedDebuffs, debuff)
			}
		} else if debuff.Duration == -1 { // 영구 디버프
			updatedDebuffs = append(updatedDebuffs, debuff)
		}
	}
	
	return updatedDebuffs
}

// GetAIInfo AI 정보 가져오기
func (m *AIManager) GetAIInfo(aiName string) (map[string]interface{}, error) {
	ai, err := m.GetAI(aiName)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"name":          ai.GetName(),
		"behavior_type": ai.GetBehaviorType(),
	}, nil
}

// GetAllAINames 등록된 모든 AI 이름 목록 반환
func (m *AIManager) GetAllAINames() []string {
	var names []string
	for name := range m.registry.ais {
		names = append(names, name)
	}
	return names
}

// RegisterCustomAI 커스텀 AI 등록
func (m *AIManager) RegisterCustomAI(name string, ai EnemyAI) {
	m.registry.Register(name, ai)
}

// ValidateAIAction AI 행동의 유효성 검사
func (m *AIManager) ValidateAIAction(
	enemyState *domain.EnemyState,
	playerState *domain.PlayerState,
	gameState *domain.GameState,
	turnNumber int,
	floorNumber int,
	aiName string,
	actionType string,
) (bool, string, error) {
	// AI 가져오기
	ai, err := m.GetAI(aiName)
	if err != nil {
		return false, "", fmt.Errorf("AI 검증 중 오류: %w", err)
	}
	
	// AI 컨텍스트 생성
	ctx := &AIContext{
		EnemyState:  enemyState,
		PlayerState: playerState,
		GameState:   gameState,
		TurnNumber:  turnNumber,
		FloorNumber: floorNumber,
	}
	
	// 행동 유효성 검사
	canExecute, reason := ai.CanExecuteAction(ctx, actionType)
	return canExecute, reason, nil
}