package rewards

import (
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// RewardManagerImpl 보상 매니저 구현
type RewardManagerImpl struct {
	generator  RewardGenerator
	repository RewardRepository
	cardRepo   domain.CardRepository
	userRepo   domain.UserRepository
}

// NewRewardManager 새로운 보상 매니저 생성
func NewRewardManager(
	generator RewardGenerator,
	repository RewardRepository,
	cardRepo domain.CardRepository,
	userRepo domain.UserRepository,
) *RewardManagerImpl {
	return &RewardManagerImpl{
		generator:  generator,
		repository: repository,
		cardRepo:   cardRepo,
		userRepo:   userRepo,
	}
}

// ProcessRewards 전투 승리 후 보상 처리
func (m *RewardManagerImpl) ProcessRewards(
	sessionID string,
	playerState *domain.PlayerState,
	gameState *domain.GameState,
	ctx *RewardContext,
) (*RewardBundle, error) {
	// 보상 생성
	bundle, err := m.generator.GenerateRewards(ctx)
	if err != nil {
		return nil, fmt.Errorf("보상 생성 실패: %w", err)
	}

	// 기본 보상 즉시 적용 (골드, 체력 등)
	for _, reward := range bundle.BaseRewards {
		err := m.ApplyReward(sessionID, playerState, gameState, &reward)
		if err != nil {
			return nil, fmt.Errorf("기본 보상 적용 실패: %w", err)
		}
	}

	// 보상 묶음 저장
	err = m.repository.SaveRewardBundle(sessionID, bundle)
	if err != nil {
		return nil, fmt.Errorf("보상 저장 실패: %w", err)
	}

	return bundle, nil
}

// ApplyReward 개별 보상 적용
func (m *RewardManagerImpl) ApplyReward(
	sessionID string,
	playerState *domain.PlayerState,
	gameState *domain.GameState,
	reward *Reward,
) error {
	switch reward.Type {
	case RewardTypeGold:
		return m.applyGoldReward(gameState, reward)
	
	case RewardTypeCard:
		return m.applyCardReward(sessionID, playerState, reward)
	
	case RewardTypeRelic:
		return m.applyRelicReward(playerState, gameState, reward)
	
	case RewardTypePotion:
		return m.applyPotionReward(gameState, reward)
	
	case RewardTypeHealth:
		return m.applyHealthReward(playerState, reward)
	
	case RewardTypeUpgrade:
		return m.applyUpgradeReward(sessionID, playerState, reward)
	
	default:
		return fmt.Errorf("지원하지 않는 보상 타입: %s", reward.Type)
	}
}

// ValidateRewardChoice 보상 선택 유효성 검사
func (m *RewardManagerImpl) ValidateRewardChoice(bundleID string, rewardID string) (bool, string) {
	// TODO: 실제 구현에서는 세션 정보를 통해 검증
	if bundleID == "" {
		return false, "보상 묶음 ID가 필요합니다"
	}
	
	if rewardID == "" {
		return false, "보상 ID가 필요합니다"
	}
	
	return true, ""
}

// CompleteRewardSelection 보상 선택 완료
func (m *RewardManagerImpl) CompleteRewardSelection(
	sessionID string,
	bundleID string,
	selectedRewardIDs []string,
	playerState *domain.PlayerState,
	gameState *domain.GameState,
) error {
	// 보상 묶음 조회
	bundle, err := m.repository.GetRewardBundle(sessionID, bundleID)
	if err != nil {
		return fmt.Errorf("보상 묶음 조회 실패: %w", err)
	}

	if bundle.IsCompleted {
		return fmt.Errorf("이미 완료된 보상입니다")
	}

	// 선택된 보상들 적용
	for _, rewardID := range selectedRewardIDs {
		reward := m.findRewardInBundle(bundle, rewardID)
		if reward == nil {
			return fmt.Errorf("보상을 찾을 수 없습니다: %s", rewardID)
		}

		err := m.ApplyReward(sessionID, playerState, gameState, reward)
		if err != nil {
			return fmt.Errorf("보상 적용 실패: %w", err)
		}
	}

	// 보상 완료 처리
	err = m.repository.MarkRewardCompleted(sessionID, bundleID)
	if err != nil {
		return fmt.Errorf("보상 완료 처리 실패: %w", err)
	}

	return nil
}

// applyGoldReward 골드 보상 적용
func (m *RewardManagerImpl) applyGoldReward(gameState *domain.GameState, reward *Reward) error {
	gameState.Gold += reward.Value
	return nil
}

// applyCardReward 카드 보상 적용
func (m *RewardManagerImpl) applyCardReward(sessionID string, playerState *domain.PlayerState, reward *Reward) error {
	// 카드를 플레이어 덱에 추가
	// TODO: 실제 구현에서는 UserCardRepository를 통해 카드 인스턴스 생성
	
	// 임시로 덱에 카드 ID 추가 (실제로는 카드 인스턴스 ID를 사용)
	if playerState.Deck == nil {
		playerState.Deck = []string{}
	}
	
	playerState.Deck = append(playerState.Deck, reward.ItemID)
	
	return nil
}

// applyRelicReward 유물 보상 적용
func (m *RewardManagerImpl) applyRelicReward(playerState *domain.PlayerState, gameState *domain.GameState, reward *Reward) error {
	// 유물을 게임 상태에 추가
	if gameState.Relics == nil {
		gameState.Relics = []string{}
	}
	
	gameState.Relics = append(gameState.Relics, reward.ItemID)
	
	// 유물 효과 즉시 적용 (예시)
	switch reward.ItemID {
	case "relic_001": // 사이버 코어
		// 최대 에너지 증가는 전투 시작시 적용되므로 여기서는 기록만
	case "relic_002": // 나노 실드
		// 전투 시작시 방어막은 전투 시작시 적용
	case "relic_003": // 데이터 크리스털
		// 드로우 확률 증가는 카드 드로우시 적용
	case "relic_004": // 양자 프로세서
		// 카드 비용 감소는 카드 플레이시 적용
	}
	
	return nil
}

// applyPotionReward 포션 보상 적용
func (m *RewardManagerImpl) applyPotionReward(gameState *domain.GameState, reward *Reward) error {
	// 포션 슬롯에 추가
	if gameState.Potions == nil {
		gameState.Potions = []string{}
	}
	
	// 포션 슬롯이 가득 차지 않았다면 추가
	if len(gameState.Potions) < gameState.PotionSlots {
		gameState.Potions = append(gameState.Potions, reward.ItemID)
	} else {
		return fmt.Errorf("포션 슬롯이 가득 참")
	}
	
	return nil
}

// applyHealthReward 체력 회복 보상 적용
func (m *RewardManagerImpl) applyHealthReward(playerState *domain.PlayerState, reward *Reward) error {
	playerState.Health += reward.Value
	
	// 최대 체력 초과 방지
	if playerState.Health > playerState.MaxHealth {
		playerState.Health = playerState.MaxHealth
	}
	
	return nil
}

// applyUpgradeReward 카드 업그레이드 보상 적용
func (m *RewardManagerImpl) applyUpgradeReward(sessionID string, playerState *domain.PlayerState, reward *Reward) error {
	// TODO: 카드 업그레이드 시스템 구현 후 연동
	// 현재는 메타데이터에 업그레이드할 카드 정보가 있다고 가정
	
	if cardID, ok := reward.Metadata["target_card_id"]; ok {
		if cardIDStr, ok := cardID.(string); ok {
			// 실제 업그레이드 로직 (추후 구현)
			_ = cardIDStr
			return nil
		}
	}
	
	return fmt.Errorf("업그레이드 대상 카드를 찾을 수 없습니다")
}

// findRewardInBundle 보상 묶음에서 특정 보상 찾기
func (m *RewardManagerImpl) findRewardInBundle(bundle *RewardBundle, rewardID string) *Reward {
	// 기본 보상에서 찾기
	for _, reward := range bundle.BaseRewards {
		if reward.ID == rewardID {
			return &reward
		}
	}
	
	// 선택 보상에서 찾기
	for _, reward := range bundle.ChoiceRewards {
		if reward.ID == rewardID {
			return &reward
		}
	}
	
	return nil
}

// GetPendingRewards 대기 중인 보상 목록 조회
func (m *RewardManagerImpl) GetPendingRewards(sessionID string) ([]*RewardBundle, error) {
	return m.repository.GetPendingRewards(sessionID)
}

// GetRewardHistory 보상 히스토리 조회
func (m *RewardManagerImpl) GetRewardHistory(sessionID string) ([]*RewardBundle, error) {
	return m.repository.GetRewardHistory(sessionID)
}

// CalculateSessionRewards 세션 전체 보상 통계
func (m *RewardManagerImpl) CalculateSessionRewards(sessionID string) (map[string]interface{}, error) {
	history, err := m.GetRewardHistory(sessionID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_gold":     0,
		"total_cards":    0,
		"total_relics":   0,
		"total_potions":  0,
		"total_healing":  0,
		"total_upgrades": 0,
	}

	for _, bundle := range history {
		// 기본 보상 집계
		for _, reward := range bundle.BaseRewards {
			m.addRewardToStats(stats, &reward)
		}
		
		// 선택 보상 집계 (완료된 것만)
		if bundle.IsCompleted {
			for _, reward := range bundle.ChoiceRewards {
				m.addRewardToStats(stats, &reward)
			}
		}
	}

	return stats, nil
}

// addRewardToStats 통계에 보상 추가
func (m *RewardManagerImpl) addRewardToStats(stats map[string]interface{}, reward *Reward) {
	switch reward.Type {
	case RewardTypeGold:
		stats["total_gold"] = stats["total_gold"].(int) + reward.Value
	case RewardTypeCard:
		stats["total_cards"] = stats["total_cards"].(int) + 1
	case RewardTypeRelic:
		stats["total_relics"] = stats["total_relics"].(int) + 1
	case RewardTypePotion:
		stats["total_potions"] = stats["total_potions"].(int) + 1
	case RewardTypeHealth:
		stats["total_healing"] = stats["total_healing"].(int) + reward.Value
	case RewardTypeUpgrade:
		stats["total_upgrades"] = stats["total_upgrades"].(int) + 1
	}
}