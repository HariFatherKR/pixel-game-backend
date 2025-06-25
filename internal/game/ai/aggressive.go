package ai

import (
	"fmt"
	"math/rand"
	"github.com/yourusername/pixel-game/internal/domain"
)

// AggressiveAI 공격적인 AI - 주로 공격 위주의 행동
type AggressiveAI struct {
	baseDamage   int
	damageScaling float64
	specialChance float64 // 특수 공격 확률
}

// NewAggressiveAI 새로운 공격적 AI 생성
func NewAggressiveAI(baseDamage int, damageScaling float64) *AggressiveAI {
	return &AggressiveAI{
		baseDamage:    baseDamage,
		damageScaling: damageScaling,
		specialChance: 0.2, // 20% 확률로 특수 공격
	}
}

// GetName AI 이름 반환
func (ai *AggressiveAI) GetName() string {
	return "Aggressive"
}

// GetBehaviorType AI 행동 유형 반환
func (ai *AggressiveAI) GetBehaviorType() string {
	return string(BehaviorAggressive)
}

// CalculateIntent 다음 턴 의도 계산
func (ai *AggressiveAI) CalculateIntent(ctx *AIContext) (*domain.EnemyIntent, error) {
	// 공격적 AI는 80% 확률로 공격, 20% 확률로 특수 행동
	if ai.shouldUseSpecialAction(ctx) {
		return ai.calculateSpecialIntent(ctx), nil
	}
	
	return ai.calculateAttackIntent(ctx), nil
}

// ExecuteAction 현재 턴 행동 실행
func (ai *AggressiveAI) ExecuteAction(ctx *AIContext) (*AIResult, error) {
	// 현재 의도에 따라 행동 실행
	switch ctx.EnemyState.Intent.Type {
	case "ATTACK":
		return ai.executeAttack(ctx)
	case "SPECIAL_ATTACK":
		return ai.executeSpecialAttack(ctx)
	case "BUFF":
		return ai.executeBuff(ctx)
	default:
		return ai.executeAttack(ctx) // 기본값은 공격
	}
}

// CanExecuteAction 행동 실행 가능 여부 검사
func (ai *AggressiveAI) CanExecuteAction(ctx *AIContext, actionType string) (bool, string) {
	switch actionType {
	case "ATTACK":
		return true, ""
	case "SPECIAL_ATTACK":
		// 체력이 50% 이하일 때만 특수 공격 가능
		if ctx.EnemyState.Health <= ctx.EnemyState.MaxHealth/2 {
			return true, ""
		}
		return false, "체력이 충분할 때는 특수 공격 불가"
	case "BUFF":
		// 이미 버프가 있으면 중복 불가
		for _, buff := range ctx.EnemyState.Buffs {
			if buff.BuffID == "rage" {
				return false, "이미 분노 버프가 활성화됨"
			}
		}
		return true, ""
	default:
		return false, "지원하지 않는 행동 타입"
	}
}

// calculateAttackIntent 공격 의도 계산
func (ai *AggressiveAI) calculateAttackIntent(ctx *AIContext) *domain.EnemyIntent {
	damage := ai.calculateDamage(ctx)
	
	return &domain.EnemyIntent{
		Type:        "ATTACK",
		Value:       damage,
		Description: fmt.Sprintf("%d 데미지 공격 준비 중", damage),
	}
}

// calculateSpecialIntent 특수 행동 의도 계산
func (ai *AggressiveAI) calculateSpecialIntent(ctx *AIContext) *domain.EnemyIntent {
	// 체력이 낮으면 분노 버프, 아니면 강력한 공격
	if ctx.EnemyState.Health <= ctx.EnemyState.MaxHealth/3 {
		return &domain.EnemyIntent{
			Type:        "BUFF",
			Value:       2,
			Description: "분노 상태로 진입 중 (공격력 +2)",
		}
	}
	
	damage := int(float64(ai.calculateDamage(ctx)) * 1.5)
	return &domain.EnemyIntent{
		Type:        "SPECIAL_ATTACK",
		Value:       damage,
		Description: fmt.Sprintf("강력한 공격 준비 중 (%d 데미지)", damage),
	}
}

// executeAttack 기본 공격 실행
func (ai *AggressiveAI) executeAttack(ctx *AIContext) (*AIResult, error) {
	damage := ai.calculateDamage(ctx)
	
	// 플레이어에게 데미지 적용
	actualDamage := ai.applyDamageToPlayer(ctx.PlayerState, damage)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "ATTACK",
			TargetID:    "player",
			Value:       damage,
			Description: fmt.Sprintf("%d 데미지로 공격", damage),
		},
		Damage:   actualDamage,
		Messages: []string{fmt.Sprintf("적이 %d 데미지로 공격했습니다!", actualDamage)},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// executeSpecialAttack 특수 공격 실행
func (ai *AggressiveAI) executeSpecialAttack(ctx *AIContext) (*AIResult, error) {
	damage := int(float64(ai.calculateDamage(ctx)) * 1.5)
	
	// 플레이어에게 데미지 적용
	actualDamage := ai.applyDamageToPlayer(ctx.PlayerState, damage)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "SPECIAL_ATTACK",
			TargetID:    "player",
			Value:       damage,
			Description: fmt.Sprintf("강력한 공격으로 %d 데미지", damage),
		},
		Damage:   actualDamage,
		Messages: []string{fmt.Sprintf("적이 강력한 공격으로 %d 데미지를 입혔습니다!", actualDamage)},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// executeBuff 버프 실행 (분노 상태)
func (ai *AggressiveAI) executeBuff(ctx *AIContext) (*AIResult, error) {
	// 분노 버프 추가
	rageBuff := domain.BuffState{
		BuffID:      "rage",
		Name:        "분노",
		Description: "공격력이 2 증가",
		Value:       2,
		Duration:    3, // 3턴 지속
	}
	
	ctx.EnemyState.Buffs = append(ctx.EnemyState.Buffs, rageBuff)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "BUFF",
			TargetID:    "self",
			Value:       2,
			Description: "분노 상태로 진입",
		},
		Buffs:    []domain.BuffState{rageBuff},
		Messages: []string{"적이 분노 상태로 진입했습니다! 공격력이 증가합니다."},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// calculateDamage 데미지 계산 (버프/디버프 고려)
func (ai *AggressiveAI) calculateDamage(ctx *AIContext) int {
	baseDamage := ai.baseDamage + int(float64(ctx.FloorNumber)*ai.damageScaling)
	
	// 분노 버프 확인
	for _, buff := range ctx.EnemyState.Buffs {
		if buff.BuffID == "rage" {
			baseDamage += buff.Value
		}
	}
	
	// 약화 디버프 확인
	for _, debuff := range ctx.EnemyState.Debuffs {
		if debuff.DebuffID == "weak" {
			baseDamage = int(float64(baseDamage) * 0.75)
		}
	}
	
	return baseDamage
}

// applyDamageToPlayer 플레이어에게 데미지 적용
func (ai *AggressiveAI) applyDamageToPlayer(player *domain.PlayerState, damage int) int {
	// 플레이어의 방어막 고려
	if player.Shield > 0 {
		if player.Shield >= damage {
			player.Shield -= damage
			return damage // 방어막이 모든 데미지를 흡수
		}
		// 방어막이 일부만 흡수
		remainingDamage := damage - player.Shield
		player.Shield = 0
		player.Health -= remainingDamage
		if player.Health < 0 {
			player.Health = 0
		}
		return damage
	}
	
	// 방어막이 없으면 직접 체력에 데미지
	player.Health -= damage
	if player.Health < 0 {
		player.Health = 0
	}
	
	return damage
}

// shouldUseSpecialAction 특수 행동을 사용할지 결정
func (ai *AggressiveAI) shouldUseSpecialAction(ctx *AIContext) bool {
	// 체력이 낮거나 랜덤 확률
	lowHealth := ctx.EnemyState.Health <= ctx.EnemyState.MaxHealth/3
	randomChance := rand.Float64() < ai.specialChance
	
	return lowHealth || randomChance
}