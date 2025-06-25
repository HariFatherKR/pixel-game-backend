package rewards

import (
	"fmt"
	"math/rand"
	"github.com/google/uuid"
	"github.com/yourusername/pixel-game/internal/domain"
)

// BasicRewardGenerator 기본 보상 생성기
type BasicRewardGenerator struct {
	cardRepo domain.CardRepository
}

// NewBasicRewardGenerator 새로운 기본 보상 생성기 생성
func NewBasicRewardGenerator(cardRepo domain.CardRepository) *BasicRewardGenerator {
	return &BasicRewardGenerator{
		cardRepo: cardRepo,
	}
}

// GenerateRewards 전투 승리 보상 묶음 생성
func (g *BasicRewardGenerator) GenerateRewards(ctx *RewardContext) (*RewardBundle, error) {
	bundle := &RewardBundle{
		ID:          uuid.New().String(),
		SourceType:  "COMBAT",
		SourceID:    fmt.Sprintf("enemy_floor_%d", ctx.FloorNumber),
		FloorNumber: ctx.FloorNumber,
		BaseRewards: []Reward{},
		ChoiceRewards: []Reward{},
		IsCompleted: false,
	}

	// 기본 골드 보상 (항상 지급)
	goldReward, err := g.GenerateGoldReward(ctx)
	if err != nil {
		return nil, fmt.Errorf("골드 보상 생성 실패: %w", err)
	}
	bundle.BaseRewards = append(bundle.BaseRewards, *goldReward)

	// 카드 보상 (선택 가능)
	cardRewards, err := g.GenerateCardRewards(ctx, 3) // 3장 중 1장 선택
	if err != nil {
		return nil, fmt.Errorf("카드 보상 생성 실패: %w", err)
	}
	bundle.ChoiceRewards = append(bundle.ChoiceRewards, cardRewards...)

	// 보스전이면 유물 보상 추가
	if ctx.FloorNumber%3 == 0 || ctx.EnemyType == "ELITE" {
		relicReward, err := g.GenerateRelicReward(ctx)
		if err == nil { // 유물 생성 실패해도 계속 진행
			bundle.ChoiceRewards = append(bundle.ChoiceRewards, *relicReward)
		}
	}

	// 랜덤 추가 보상 (낮은 확률)
	if rand.Float64() < 0.15 { // 15% 확률
		extraReward := g.generateExtraReward(ctx)
		if extraReward != nil {
			bundle.BaseRewards = append(bundle.BaseRewards, *extraReward)
		}
	}

	return bundle, nil
}

// GenerateCardRewards 카드 보상 생성
func (g *BasicRewardGenerator) GenerateCardRewards(ctx *RewardContext, count int) ([]Reward, error) {
	rewards := []Reward{}

	// 카드 등급별 확률 계산
	rarityWeights := g.calculateCardRarityWeights(ctx)

	for i := 0; i < count; i++ {
		// 등급 결정
		rarity := g.selectRarityByWeight(rarityWeights)
		
		// 해당 등급의 카드 목록 가져오기
		cardFilter := domain.CardFilter{
			Rarity: (*domain.CardRarity)(&rarity),
			Limit:  20,
		}
		
		cards, err := g.cardRepo.GetAll(cardFilter)
		if err != nil || len(cards) == 0 {
			// 실패시 일반 등급으로 대체
			rarity = RewardRarityCommon
			cardFilter.Rarity = (*domain.CardRarity)(&rarity)
			cards, err = g.cardRepo.GetAll(cardFilter)
			if err != nil || len(cards) == 0 {
				continue // 이 카드는 건너뛰기
			}
		}

		// 랜덤 카드 선택
		selectedCard := cards[rand.Intn(len(cards))]

		reward := Reward{
			ID:          uuid.New().String(),
			Type:        RewardTypeCard,
			Rarity:      rarity,
			ItemID:      selectedCard.ID,
			Name:        selectedCard.Name,
			Description: selectedCard.Description,
			ImageURL:    selectedCard.ImageURL,
			Metadata: map[string]interface{}{
				"card_type":   selectedCard.Type,
				"card_cost":   selectedCard.Cost,
				"card_rarity": selectedCard.Rarity,
			},
		}

		rewards = append(rewards, reward)
	}

	return rewards, nil
}

// GenerateGoldReward 골드 보상 생성
func (g *BasicRewardGenerator) GenerateGoldReward(ctx *RewardContext) (*Reward, error) {
	baseGold := g.CalculateRewardValue(RewardTypeGold, ctx)
	
	// ±20% 랜덤 변동
	variation := int(float64(baseGold) * 0.2)
	finalGold := baseGold + rand.Intn(variation*2+1) - variation

	if finalGold < 1 {
		finalGold = 1
	}

	reward := &Reward{
		ID:          uuid.New().String(),
		Type:        RewardTypeGold,
		Rarity:      RewardRarityCommon,
		Value:       finalGold,
		Name:        fmt.Sprintf("골드 %d", finalGold),
		Description: fmt.Sprintf("%d 골드를 획득합니다", finalGold),
		ImageURL:    "/images/rewards/gold.png",
		Metadata: map[string]interface{}{
			"base_amount": baseGold,
			"floor":       ctx.FloorNumber,
		},
	}

	return reward, nil
}

// GenerateRelicReward 유물 보상 생성
func (g *BasicRewardGenerator) GenerateRelicReward(ctx *RewardContext) (*Reward, error) {
	// 유물 목록 (향후 별도 레포지토리로 분리 예정)
	relics := []struct {
		ID          string
		Name        string
		Description string
		Rarity      RewardRarity
	}{
		{"relic_001", "사이버 코어", "매 턴 시작시 에너지 +1", RewardRarityRare},
		{"relic_002", "나노 실드", "전투 시작시 방어막 +5", RewardRarityRare},
		{"relic_003", "데이터 크리스털", "카드 드로우 시 25% 확률로 추가 드로우", RewardRarityEpic},
		{"relic_004", "양자 프로세서", "카드 비용 1 감소 (최소 0)", RewardRarityLegendary},
	}

	// 층수에 따른 등급 가중치
	var selectedRelic *struct {
		ID          string
		Name        string
		Description string
		Rarity      RewardRarity
	}

	rarityWeights := map[RewardRarity]float64{
		RewardRarityRare:      0.6,
		RewardRarityEpic:      0.3,
		RewardRarityLegendary: 0.1,
	}

	// 층수가 높을수록 좋은 유물 확률 증가
	if ctx.FloorNumber >= 6 {
		rarityWeights[RewardRarityEpic] = 0.5
		rarityWeights[RewardRarityLegendary] = 0.2
		rarityWeights[RewardRarityRare] = 0.3
	}

	targetRarity := g.selectRarityByWeight(rarityWeights)

	// 해당 등급의 유물 필터링
	availableRelics := []struct {
		ID          string
		Name        string
		Description string
		Rarity      RewardRarity
	}{}

	for _, relic := range relics {
		if relic.Rarity == targetRarity {
			availableRelics = append(availableRelics, relic)
		}
	}

	if len(availableRelics) == 0 {
		// 대체 유물 선택
		selectedRelic = &relics[0]
	} else {
		selectedRelic = &availableRelics[rand.Intn(len(availableRelics))]
	}

	reward := &Reward{
		ID:          uuid.New().String(),
		Type:        RewardTypeRelic,
		Rarity:      selectedRelic.Rarity,
		ItemID:      selectedRelic.ID,
		Name:        selectedRelic.Name,
		Description: selectedRelic.Description,
		ImageURL:    fmt.Sprintf("/images/relics/%s.png", selectedRelic.ID),
		Metadata: map[string]interface{}{
			"relic_type": "passive",
			"floor":      ctx.FloorNumber,
		},
	}

	return reward, nil
}

// CalculateRewardValue 보상 가치 계산
func (g *BasicRewardGenerator) CalculateRewardValue(rewardType RewardType, ctx *RewardContext) int {
	switch rewardType {
	case RewardTypeGold:
		// 기본 골드 = 20 + (층수 * 5) + 난이도 보정
		baseGold := 20 + (ctx.FloorNumber * 5)
		difficultyBonus := int(float64(baseGold) * ctx.DifficultyMod)
		return baseGold + difficultyBonus

	case RewardTypeHealth:
		// 체력 회복량 = 최대 체력의 25% + 층수 보정
		if ctx.PlayerState != nil {
			baseHeal := ctx.PlayerState.MaxHealth / 4
			floorBonus := ctx.FloorNumber * 2
			return baseHeal + floorBonus
		}
		return 15 // 기본값

	default:
		return 0
	}
}

// calculateCardRarityWeights 카드 등급별 가중치 계산
func (g *BasicRewardGenerator) calculateCardRarityWeights(ctx *RewardContext) map[RewardRarity]float64 {
	weights := map[RewardRarity]float64{
		RewardRarityCommon:    0.6,
		RewardRarityRare:      0.3,
		RewardRarityEpic:      0.09,
		RewardRarityLegendary: 0.01,
	}

	// 층수가 높을수록 좋은 카드 확률 증가
	floorMod := float64(ctx.FloorNumber) * 0.02
	if floorMod > 0.3 {
		floorMod = 0.3
	}

	weights[RewardRarityCommon] -= floorMod
	weights[RewardRarityRare] += floorMod * 0.6
	weights[RewardRarityEpic] += floorMod * 0.3
	weights[RewardRarityLegendary] += floorMod * 0.1

	// 음수 방지
	for rarity := range weights {
		if weights[rarity] < 0.01 {
			weights[rarity] = 0.01
		}
	}

	return weights
}

// selectRarityByWeight 가중치에 따른 등급 선택
func (g *BasicRewardGenerator) selectRarityByWeight(weights map[RewardRarity]float64) RewardRarity {
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}

	randomValue := rand.Float64() * totalWeight
	currentWeight := 0.0

	for rarity, weight := range weights {
		currentWeight += weight
		if randomValue <= currentWeight {
			return rarity
		}
	}

	return RewardRarityCommon // 기본값
}

// generateExtraReward 추가 보상 생성 (포션, 체력 회복 등)
func (g *BasicRewardGenerator) generateExtraReward(ctx *RewardContext) *Reward {
	extraTypes := []RewardType{RewardTypePotion, RewardTypeHealth}
	selectedType := extraTypes[rand.Intn(len(extraTypes))]

	switch selectedType {
	case RewardTypePotion:
		return &Reward{
			ID:          uuid.New().String(),
			Type:        RewardTypePotion,
			Rarity:      RewardRarityCommon,
			ItemID:      "potion_heal",
			Name:        "치유 포션",
			Description: "체력을 25 회복합니다",
			ImageURL:    "/images/potions/heal.png",
			Value:       25,
		}

	case RewardTypeHealth:
		healAmount := g.CalculateRewardValue(RewardTypeHealth, ctx)
		return &Reward{
			ID:          uuid.New().String(),
			Type:        RewardTypeHealth,
			Rarity:      RewardRarityCommon,
			Value:       healAmount,
			Name:        fmt.Sprintf("체력 회복 %d", healAmount),
			Description: fmt.Sprintf("%d 체력을 즉시 회복합니다", healAmount),
			ImageURL:    "/images/rewards/health.png",
		}

	default:
		return nil
	}
}