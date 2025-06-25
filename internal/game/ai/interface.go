package ai

import (
	"github.com/yourusername/pixel-game/internal/domain"
)

// AIContext 적 AI가 결정을 내리기 위한 게임 상황 정보
type AIContext struct {
	EnemyState  *domain.EnemyState
	PlayerState *domain.PlayerState
	GameState   *domain.GameState
	TurnNumber  int
	FloorNumber int
}

// AIAction 적이 수행할 수 있는 행동
type AIAction struct {
	Type        string                 `json:"type"`        // "ATTACK", "DEFEND", "BUFF", "DEBUFF", "SPECIAL"
	TargetID    string                 `json:"target_id"`   // 대상 ID (보통 플레이어)
	Value       int                    `json:"value"`       // 행동의 수치값 (데미지, 방어력 등)
	Description string                 `json:"description"` // 행동 설명
	Parameters  map[string]interface{} `json:"parameters"`  // 추가 매개변수
}

// AIResult AI 행동 실행 결과
type AIResult struct {
	Success     bool           `json:"success"`
	Action      AIAction       `json:"action"`
	Damage      int            `json:"damage,omitempty"`
	Shield      int            `json:"shield,omitempty"`
	Buffs       []domain.BuffState   `json:"buffs,omitempty"`
	Debuffs     []domain.DebuffState `json:"debuffs,omitempty"`
	Messages    []string       `json:"messages"`
	NextIntent  *domain.EnemyIntent `json:"next_intent,omitempty"`
}

// EnemyAI 적 AI의 기본 인터페이스
type EnemyAI interface {
	// GetName AI의 이름을 반환
	GetName() string
	
	// GetBehaviorType AI의 행동 유형을 반환 (AGGRESSIVE, DEFENSIVE, BALANCED)
	GetBehaviorType() string
	
	// CalculateIntent 현재 상황에서 다음 턴 의도를 계산
	CalculateIntent(ctx *AIContext) (*domain.EnemyIntent, error)
	
	// ExecuteAction 현재 턴에 행동을 실행
	ExecuteAction(ctx *AIContext) (*AIResult, error)
	
	// CanExecuteAction 해당 행동을 실행할 수 있는지 검사
	CanExecuteAction(ctx *AIContext, actionType string) (bool, string)
}

// AIBehaviorType AI 행동 유형 상수
type AIBehaviorType string

const (
	BehaviorAggressive AIBehaviorType = "AGGRESSIVE" // 공격적 - 주로 공격 위주
	BehaviorDefensive  AIBehaviorType = "DEFENSIVE"  // 방어적 - 방어와 회복 위주
	BehaviorBalanced   AIBehaviorType = "BALANCED"   // 균형 - 상황에 따라 적절히 선택
	BehaviorSpecial    AIBehaviorType = "SPECIAL"    // 특수 - 고유한 패턴
)

// AIActionType AI가 수행할 수 있는 행동 타입
type AIActionType string

const (
	ActionAttack  AIActionType = "ATTACK"  // 공격
	ActionDefend  AIActionType = "DEFEND"  // 방어
	ActionBuff    AIActionType = "BUFF"    // 자신에게 버프
	ActionDebuff  AIActionType = "DEBUFF"  // 플레이어에게 디버프
	ActionSpecial AIActionType = "SPECIAL" // 특수 능력
	ActionHeal    AIActionType = "HEAL"    // 회복
)

// AIRegistry AI들을 관리하는 레지스트리
type AIRegistry struct {
	ais map[string]EnemyAI
}

// NewAIRegistry 새로운 AI 레지스트리 생성
func NewAIRegistry() *AIRegistry {
	return &AIRegistry{
		ais: make(map[string]EnemyAI),
	}
}

// Register AI를 레지스트리에 등록
func (r *AIRegistry) Register(name string, ai EnemyAI) {
	r.ais[name] = ai
}

// Get 이름으로 AI를 가져옴
func (r *AIRegistry) Get(name string) (EnemyAI, bool) {
	ai, exists := r.ais[name]
	return ai, exists
}

// GetAll 모든 등록된 AI 목록을 반환
func (r *AIRegistry) GetAll() map[string]EnemyAI {
	return r.ais
}