package ai

import (
	"fmt"
	"math/rand"
	"github.com/yourusername/pixel-game/internal/domain"
)

// DefensiveAI 방어적인 AI - 방어와 회복 위주의 행동
type DefensiveAI struct {
	baseDamage    int
	baseShield    int
	healAmount    int
	defensiveMode bool // 방어 모드 여부
}

// NewDefensiveAI 새로운 방어적 AI 생성
func NewDefensiveAI(baseDamage, baseShield, healAmount int) *DefensiveAI {
	return &DefensiveAI{
		baseDamage:    baseDamage,
		baseShield:    baseShield,
		healAmount:    healAmount,
		defensiveMode: false,
	}
}

// GetName AI 이름 반환
func (ai *DefensiveAI) GetName() string {
	return "Defensive"
}

// GetBehaviorType AI 행동 유형 반환
func (ai *DefensiveAI) GetBehaviorType() string {
	return string(BehaviorDefensive)
}

// CalculateIntent 다음 턴 의도 계산
func (ai *DefensiveAI) CalculateIntent(ctx *AIContext) (*domain.EnemyIntent, error) {
	// 체력 상태에 따른 우선순위 결정
	healthPercentage := float64(ctx.EnemyState.Health) / float64(ctx.EnemyState.MaxHealth)
	
	// 체력이 30% 이하면 회복 우선
	if healthPercentage <= 0.3 && ai.canHeal(ctx) {
		return ai.calculateHealIntent(ctx), nil
	}
	
	// 체력이 50% 이하면 방어 우선
	if healthPercentage <= 0.5 {
		return ai.calculateDefendIntent(ctx), nil
	}
	
	// 플레이어가 강해 보이면 디버프 사용
	if ai.shouldDebuffPlayer(ctx) {
		return ai.calculateDebuffIntent(ctx), nil
	}
	
	// 기본적으로는 약간의 공격
	return ai.calculateAttackIntent(ctx), nil
}

// ExecuteAction 현재 턴 행동 실행
func (ai *DefensiveAI) ExecuteAction(ctx *AIContext) (*AIResult, error) {
	switch ctx.EnemyState.Intent.Type {
	case "DEFEND":
		return ai.executeDefend(ctx)
	case "HEAL":
		return ai.executeHeal(ctx)
	case "DEBUFF":
		return ai.executeDebuff(ctx)
	case "ATTACK":
		return ai.executeAttack(ctx)
	default:
		return ai.executeDefend(ctx) // 기본값은 방어
	}
}

// CanExecuteAction 행동 실행 가능 여부 검사
func (ai *DefensiveAI) CanExecuteAction(ctx *AIContext, actionType string) (bool, string) {
	switch actionType {
	case "DEFEND":
		return true, ""
	case "HEAL":
		if ctx.EnemyState.Health >= ctx.EnemyState.MaxHealth {
			return false, "이미 최대 체력입니다"
		}
		return true, ""
	case "DEBUFF":
		// 플레이어에게 이미 디버프가 있는지 확인
		for _, debuff := range ctx.PlayerState.Debuffs {
			if debuff.DebuffID == "weak" || debuff.DebuffID == "frail" {
				return false, "플레이어에게 이미 디버프가 적용됨"
			}
		}
		return true, ""
	case "ATTACK":
		return true, ""
	default:
		return false, "지원하지 않는 행동 타입"
	}
}

// calculateDefendIntent 방어 의도 계산
func (ai *DefensiveAI) calculateDefendIntent(ctx *AIContext) *domain.EnemyIntent {
	shield := ai.calculateShield(ctx)
	
	return &domain.EnemyIntent{
		Type:        "DEFEND",
		Value:       shield,
		Description: fmt.Sprintf("%d 방어막 생성 준비 중", shield),
	}
}

// calculateHealIntent 회복 의도 계산
func (ai *DefensiveAI) calculateHealIntent(ctx *AIContext) *domain.EnemyIntent {
	healAmount := ai.healAmount
	
	return &domain.EnemyIntent{
		Type:        "HEAL",
		Value:       healAmount,
		Description: fmt.Sprintf("%d 체력 회복 준비 중", healAmount),
	}
}

// calculateDebuffIntent 디버프 의도 계산
func (ai *DefensiveAI) calculateDebuffIntent(ctx *AIContext) *domain.EnemyIntent {
	return &domain.EnemyIntent{
		Type:        "DEBUFF",
		Value:       2, // 디버프 지속시간
		Description: "플레이어 약화 준비 중",
	}
}

// calculateAttackIntent 공격 의도 계산
func (ai *DefensiveAI) calculateAttackIntent(ctx *AIContext) *domain.EnemyIntent {
	damage := ai.calculateDamage(ctx)
	
	return &domain.EnemyIntent{
		Type:        "ATTACK",
		Value:       damage,
		Description: fmt.Sprintf("%d 데미지 공격 준비 중", damage),
	}
}

// executeDefend 방어 실행
func (ai *DefensiveAI) executeDefend(ctx *AIContext) (*AIResult, error) {
	shield := ai.calculateShield(ctx)
	
	// 방어막 적용
	ctx.EnemyState.Shield += shield
	ai.defensiveMode = true
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "DEFEND",
			TargetID:    "self",
			Value:       shield,
			Description: fmt.Sprintf("%d 방어막 생성", shield),
		},
		Shield:   shield,
		Messages: []string{fmt.Sprintf("적이 %d 방어막을 생성했습니다!", shield)},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// executeHeal 회복 실행
func (ai *DefensiveAI) executeHeal(ctx *AIContext) (*AIResult, error) {
	healAmount := ai.healAmount
	
	// 최대 체력을 초과하지 않도록 조정
	maxHeal := ctx.EnemyState.MaxHealth - ctx.EnemyState.Health
	if healAmount > maxHeal {
		healAmount = maxHeal
	}
	
	ctx.EnemyState.Health += healAmount
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "HEAL",
			TargetID:    "self",
			Value:       healAmount,
			Description: fmt.Sprintf("%d 체력 회복", healAmount),
		},
		Messages: []string{fmt.Sprintf("적이 %d 체력을 회복했습니다!", healAmount)},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// executeDebuff 디버프 실행
func (ai *DefensiveAI) executeDebuff(ctx *AIContext) (*AIResult, error) {
	// 랜덤하게 약화 또는 연약 디버프 적용
	var debuff domain.DebuffState
	if rand.Float64() < 0.5 {
		debuff = domain.DebuffState{
			DebuffID:    "weak",
			Name:        "약화",
			Description: "공격력 25% 감소",
			Value:       25,
			Duration:    2,
		}
	} else {
		debuff = domain.DebuffState{
			DebuffID:    "frail",
			Name:        "연약",
			Description: "방어막 획득량 25% 감소",
			Value:       25,
			Duration:    2,
		}
	}
	
	ctx.PlayerState.Debuffs = append(ctx.PlayerState.Debuffs, debuff)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "DEBUFF",
			TargetID:    "player",
			Value:       debuff.Value,
			Description: fmt.Sprintf("플레이어에게 %s 적용", debuff.Name),
		},
		Debuffs:  []domain.DebuffState{debuff},
		Messages: []string{fmt.Sprintf("적이 당신에게 %s을(를) 적용했습니다!", debuff.Name)},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// executeAttack 공격 실행
func (ai *DefensiveAI) executeAttack(ctx *AIContext) (*AIResult, error) {
	damage := ai.calculateDamage(ctx)
	
	// 플레이어에게 데미지 적용
	actualDamage := ai.applyDamageToPlayer(ctx.PlayerState, damage)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "ATTACK",
			TargetID:    "player",
			Value:       damage,
			Description: fmt.Sprintf("%d 데미지 공격", damage),
		},
		Damage:   actualDamage,
		Messages: []string{fmt.Sprintf("적이 %d 데미지로 공격했습니다!", actualDamage)},
	}
	
	// 다음 의도 계산
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	
	return result, nil
}

// calculateDamage 데미지 계산
func (ai *DefensiveAI) calculateDamage(ctx *AIContext) int {
	// 방어적 AI는 공격력이 낮음
	damage := int(float64(ai.baseDamage) * 0.8)
	
	// 약화 디버프 확인
	for _, debuff := range ctx.EnemyState.Debuffs {
		if debuff.DebuffID == "weak" {
			damage = int(float64(damage) * 0.75)
		}
	}
	
	return damage
}

// calculateShield 방어막 계산
func (ai *DefensiveAI) calculateShield(ctx *AIContext) int {
	shield := ai.baseShield
	
	// 방어 모드일 때 방어막 증가
	if ai.defensiveMode {
		shield = int(float64(shield) * 1.3)
	}
	
	return shield
}

// applyDamageToPlayer 플레이어에게 데미지 적용 (공격적 AI와 동일)
func (ai *DefensiveAI) applyDamageToPlayer(player *domain.PlayerState, damage int) int {
	if player.Shield > 0 {
		if player.Shield >= damage {
			player.Shield -= damage
			return damage
		}
		remainingDamage := damage - player.Shield
		player.Shield = 0
		player.Health -= remainingDamage
		if player.Health < 0 {
			player.Health = 0
		}
		return damage
	}
	
	player.Health -= damage
	if player.Health < 0 {
		player.Health = 0
	}
	
	return damage
}

// canHeal 회복 가능 여부 확인
func (ai *DefensiveAI) canHeal(ctx *AIContext) bool {
	return ctx.EnemyState.Health < ctx.EnemyState.MaxHealth
}

// shouldDebuffPlayer 플레이어에게 디버프를 적용해야 하는지 판단
func (ai *DefensiveAI) shouldDebuffPlayer(ctx *AIContext) bool {
	// 플레이어의 체력이 높고 아직 디버프가 없을 때
	playerHealthy := ctx.PlayerState.Health > ctx.PlayerState.MaxHealth/2
	noDebuffs := len(ctx.PlayerState.Debuffs) == 0
	
	return playerHealthy && noDebuffs && rand.Float64() < 0.3 // 30% 확률
}