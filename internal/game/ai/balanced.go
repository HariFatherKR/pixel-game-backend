package ai

import (
	"fmt"
	"math/rand"
	"github.com/yourusername/pixel-game/internal/domain"
)

// BalancedAI 균형잡힌 AI - 상황에 따라 적절한 행동 선택
type BalancedAI struct {
	baseDamage      int
	baseShield      int
	healAmount      int
	adaptiveCounter int // 적응적 카운터
}

// NewBalancedAI 새로운 균형잡힌 AI 생성
func NewBalancedAI(baseDamage, baseShield, healAmount int) *BalancedAI {
	return &BalancedAI{
		baseDamage:      baseDamage,
		baseShield:      baseShield,
		healAmount:      healAmount,
		adaptiveCounter: 0,
	}
}

// GetName AI 이름 반환
func (ai *BalancedAI) GetName() string {
	return "Balanced"
}

// GetBehaviorType AI 행동 유형 반환
func (ai *BalancedAI) GetBehaviorType() string {
	return string(BehaviorBalanced)
}

// CalculateIntent 다음 턴 의도 계산
func (ai *BalancedAI) CalculateIntent(ctx *AIContext) (*domain.EnemyIntent, error) {
	// 상황 분석
	situation := ai.analyzeSituation(ctx)
	
	switch situation {
	case "CRITICAL": // 위험 상황 - 회복 또는 방어
		if ai.canHeal(ctx) && ctx.EnemyState.Health <= ctx.EnemyState.MaxHealth/4 {
			return ai.calculateHealIntent(ctx), nil
		}
		return ai.calculateDefendIntent(ctx), nil
		
	case "DEFENSIVE": // 방어 상황 - 방어 또는 디버프
		if rand.Float64() < 0.6 {
			return ai.calculateDefendIntent(ctx), nil
		}
		return ai.calculateDebuffIntent(ctx), nil
		
	case "AGGRESSIVE": // 공격 상황 - 공격 또는 버프
		if ai.shouldBuff(ctx) {
			return ai.calculateBuffIntent(ctx), nil
		}
		return ai.calculateAttackIntent(ctx), nil
		
	case "NEUTRAL": // 중립 상황 - 랜덤 선택
		return ai.calculateRandomIntent(ctx), nil
		
	default:
		return ai.calculateAttackIntent(ctx), nil
	}
}

// ExecuteAction 현재 턴 행동 실행
func (ai *BalancedAI) ExecuteAction(ctx *AIContext) (*AIResult, error) {
	ai.adaptiveCounter++
	
	switch ctx.EnemyState.Intent.Type {
	case "ATTACK":
		return ai.executeAttack(ctx)
	case "DEFEND":
		return ai.executeDefend(ctx)
	case "HEAL":
		return ai.executeHeal(ctx)
	case "BUFF":
		return ai.executeBuff(ctx)
	case "DEBUFF":
		return ai.executeDebuff(ctx)
	default:
		return ai.executeAttack(ctx)
	}
}

// CanExecuteAction 행동 실행 가능 여부 검사
func (ai *BalancedAI) CanExecuteAction(ctx *AIContext, actionType string) (bool, string) {
	switch actionType {
	case "ATTACK", "DEFEND":
		return true, ""
	case "HEAL":
		if ctx.EnemyState.Health >= ctx.EnemyState.MaxHealth {
			return false, "이미 최대 체력입니다"
		}
		return true, ""
	case "BUFF":
		// 이미 강화 버프가 있는지 확인
		for _, buff := range ctx.EnemyState.Buffs {
			if buff.BuffID == "strength" || buff.BuffID == "dexterity" {
				return false, "이미 강화 버프가 활성화됨"
			}
		}
		return true, ""
	case "DEBUFF":
		// 플레이어에게 이미 디버프가 있는지 확인
		for _, debuff := range ctx.PlayerState.Debuffs {
			if debuff.DebuffID == "weak" || debuff.DebuffID == "vulnerable" {
				return false, "플레이어에게 이미 디버프가 적용됨"
			}
		}
		return true, ""
	default:
		return false, "지원하지 않는 행동 타입"
	}
}

// analyzeSituation 현재 상황 분석
func (ai *BalancedAI) analyzeSituation(ctx *AIContext) string {
	enemyHealthPercent := float64(ctx.EnemyState.Health) / float64(ctx.EnemyState.MaxHealth)
	playerHealthPercent := float64(ctx.PlayerState.Health) / float64(ctx.PlayerState.MaxHealth)
	
	// 위험 상황: 자신의 체력이 25% 이하
	if enemyHealthPercent <= 0.25 {
		return "CRITICAL"
	}
	
	// 방어 상황: 자신의 체력이 50% 이하이거나 플레이어가 강함
	if enemyHealthPercent <= 0.5 || (playerHealthPercent > 0.8 && ctx.PlayerState.Shield > 0) {
		return "DEFENSIVE"
	}
	
	// 공격 상황: 플레이어의 체력이 50% 이하이고 자신이 건강함
	if playerHealthPercent <= 0.5 && enemyHealthPercent > 0.6 {
		return "AGGRESSIVE"
	}
	
	// 중립 상황
	return "NEUTRAL"
}

// calculateAttackIntent 공격 의도 계산
func (ai *BalancedAI) calculateAttackIntent(ctx *AIContext) *domain.EnemyIntent {
	damage := ai.calculateDamage(ctx)
	
	return &domain.EnemyIntent{
		Type:        "ATTACK",
		Value:       damage,
		Description: fmt.Sprintf("%d 데미지 공격 준비 중", damage),
	}
}

// calculateDefendIntent 방어 의도 계산
func (ai *BalancedAI) calculateDefendIntent(ctx *AIContext) *domain.EnemyIntent {
	shield := ai.calculateShield(ctx)
	
	return &domain.EnemyIntent{
		Type:        "DEFEND",
		Value:       shield,
		Description: fmt.Sprintf("%d 방어막 생성 준비 중", shield),
	}
}

// calculateHealIntent 회복 의도 계산
func (ai *BalancedAI) calculateHealIntent(ctx *AIContext) *domain.EnemyIntent {
	return &domain.EnemyIntent{
		Type:        "HEAL",
		Value:       ai.healAmount,
		Description: fmt.Sprintf("%d 체력 회복 준비 중", ai.healAmount),
	}
}

// calculateBuffIntent 버프 의도 계산
func (ai *BalancedAI) calculateBuffIntent(ctx *AIContext) *domain.EnemyIntent {
	return &domain.EnemyIntent{
		Type:        "BUFF",
		Value:       2,
		Description: "자신 강화 준비 중 (공격력 +2)",
	}
}

// calculateDebuffIntent 디버프 의도 계산
func (ai *BalancedAI) calculateDebuffIntent(ctx *AIContext) *domain.EnemyIntent {
	return &domain.EnemyIntent{
		Type:        "DEBUFF",
		Value:       2,
		Description: "플레이어 약화 준비 중",
	}
}

// calculateRandomIntent 랜덤 의도 계산
func (ai *BalancedAI) calculateRandomIntent(ctx *AIContext) *domain.EnemyIntent {
	actions := []string{"ATTACK", "DEFEND"}
	
	// 상황에 따라 추가 행동 옵션
	if ai.canHeal(ctx) {
		actions = append(actions, "HEAL")
	}
	if ai.shouldBuff(ctx) {
		actions = append(actions, "BUFF")
	}
	if ai.shouldDebuff(ctx) {
		actions = append(actions, "DEBUFF")
	}
	
	selectedAction := actions[rand.Intn(len(actions))]
	
	switch selectedAction {
	case "ATTACK":
		return ai.calculateAttackIntent(ctx)
	case "DEFEND":
		return ai.calculateDefendIntent(ctx)
	case "HEAL":
		return ai.calculateHealIntent(ctx)
	case "BUFF":
		return ai.calculateBuffIntent(ctx)
	case "DEBUFF":
		return ai.calculateDebuffIntent(ctx)
	default:
		return ai.calculateAttackIntent(ctx)
	}
}

// executeAttack 공격 실행
func (ai *BalancedAI) executeAttack(ctx *AIContext) (*AIResult, error) {
	damage := ai.calculateDamage(ctx)
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
	
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	return result, nil
}

// executeDefend 방어 실행
func (ai *BalancedAI) executeDefend(ctx *AIContext) (*AIResult, error) {
	shield := ai.calculateShield(ctx)
	ctx.EnemyState.Shield += shield
	
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
	
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	return result, nil
}

// executeHeal 회복 실행
func (ai *BalancedAI) executeHeal(ctx *AIContext) (*AIResult, error) {
	healAmount := ai.healAmount
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
	
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	return result, nil
}

// executeBuff 버프 실행
func (ai *BalancedAI) executeBuff(ctx *AIContext) (*AIResult, error) {
	buff := domain.BuffState{
		BuffID:      "strength",
		Name:        "힘",
		Description: "공격력이 2 증가",
		Value:       2,
		Duration:    3,
	}
	
	ctx.EnemyState.Buffs = append(ctx.EnemyState.Buffs, buff)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "BUFF",
			TargetID:    "self",
			Value:       2,
			Description: "자신 강화",
		},
		Buffs:    []domain.BuffState{buff},
		Messages: []string{"적이 자신을 강화했습니다! 공격력이 증가합니다."},
	}
	
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	return result, nil
}

// executeDebuff 디버프 실행
func (ai *BalancedAI) executeDebuff(ctx *AIContext) (*AIResult, error) {
	debuff := domain.DebuffState{
		DebuffID:    "vulnerable",
		Name:        "취약",
		Description: "받는 데미지 50% 증가",
		Value:       50,
		Duration:    2,
	}
	
	ctx.PlayerState.Debuffs = append(ctx.PlayerState.Debuffs, debuff)
	
	result := &AIResult{
		Success: true,
		Action: AIAction{
			Type:        "DEBUFF",
			TargetID:    "player",
			Value:       50,
			Description: "플레이어에게 취약 적용",
		},
		Debuffs:  []domain.DebuffState{debuff},
		Messages: []string{"적이 당신을 취약하게 만들었습니다! 받는 데미지가 증가합니다."},
	}
	
	nextIntent, _ := ai.CalculateIntent(ctx)
	result.NextIntent = nextIntent
	return result, nil
}

// Helper 함수들
func (ai *BalancedAI) calculateDamage(ctx *AIContext) int {
	damage := ai.baseDamage
	
	// 버프 확인
	for _, buff := range ctx.EnemyState.Buffs {
		if buff.BuffID == "strength" {
			damage += buff.Value
		}
	}
	
	// 디버프 확인
	for _, debuff := range ctx.EnemyState.Debuffs {
		if debuff.DebuffID == "weak" {
			damage = int(float64(damage) * 0.75)
		}
	}
	
	return damage
}

func (ai *BalancedAI) calculateShield(ctx *AIContext) int {
	return ai.baseShield
}

func (ai *BalancedAI) applyDamageToPlayer(player *domain.PlayerState, damage int) int {
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

func (ai *BalancedAI) canHeal(ctx *AIContext) bool {
	return ctx.EnemyState.Health < ctx.EnemyState.MaxHealth
}

func (ai *BalancedAI) shouldBuff(ctx *AIContext) bool {
	// 버프가 없고 상황이 좋을 때
	for _, buff := range ctx.EnemyState.Buffs {
		if buff.BuffID == "strength" {
			return false
		}
	}
	return rand.Float64() < 0.3
}

func (ai *BalancedAI) shouldDebuff(ctx *AIContext) bool {
	// 플레이어에게 디버프가 없을 때
	for _, debuff := range ctx.PlayerState.Debuffs {
		if debuff.DebuffID == "vulnerable" || debuff.DebuffID == "weak" {
			return false
		}
	}
	return rand.Float64() < 0.4
}