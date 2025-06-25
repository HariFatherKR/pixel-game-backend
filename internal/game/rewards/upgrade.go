package rewards

import (
	"fmt"
	"github.com/yourusername/pixel-game/internal/domain"
)

// CardUpgradeServiceImpl 카드 업그레이드 서비스 구현
type CardUpgradeServiceImpl struct {
	cardRepo     domain.CardRepository
	userCardRepo domain.CardRepository // CardRepository에 UserCard 기능이 포함됨
}

// NewCardUpgradeService 새로운 카드 업그레이드 서비스 생성
func NewCardUpgradeService(
	cardRepo domain.CardRepository,
	userCardRepo domain.CardRepository,
) *CardUpgradeServiceImpl {
	return &CardUpgradeServiceImpl{
		cardRepo:     cardRepo,
		userCardRepo: userCardRepo,
	}
}

// GetUpgradeableCards 업그레이드 가능한 카드 목록
func (s *CardUpgradeServiceImpl) GetUpgradeableCards(playerState *domain.PlayerState) ([]string, error) {
	upgradeableCards := []string{}
	
	// 덱의 모든 카드 확인
	if playerState.Deck != nil {
		cardCounts := make(map[string]int)
		
		// 카드 개수 세기
		for _, cardID := range playerState.Deck {
			cardCounts[cardID]++
		}
		
		// 중복된 카드만 업그레이드 가능
		for cardID, count := range cardCounts {
			if count >= 2 {
				// 카드 정보 확인
				card, err := s.cardRepo.GetByID(cardID)
				if err != nil {
					continue // 카드를 찾을 수 없으면 건너뛰기
				}
				
				// 업그레이드 가능한 카드인지 확인
				if s.canUpgradeCardType(card) {
					upgradeableCards = append(upgradeableCards, cardID)
				}
			}
		}
	}
	
	return upgradeableCards, nil
}

// UpgradeCard 카드 업그레이드
func (s *CardUpgradeServiceImpl) UpgradeCard(cardInstanceID string, playerState *domain.PlayerState) error {
	// 카드가 덱에 있는지 확인
	cardIndex := -1
	for i, deckCardID := range playerState.Deck {
		if deckCardID == cardInstanceID {
			cardIndex = i
			break
		}
	}
	
	if cardIndex == -1 {
		return fmt.Errorf("카드가 덱에 없습니다")
	}
	
	// 원본 카드 정보 가져오기
	originalCard, err := s.cardRepo.GetByID(cardInstanceID)
	if err != nil {
		return fmt.Errorf("카드 정보 조회 실패: %w", err)
	}
	
	// 업그레이드 가능 여부 확인
	canUpgrade, reason := s.CanUpgradeCard(cardInstanceID, playerState)
	if !canUpgrade {
		return fmt.Errorf("카드 업그레이드 불가: %s", reason)
	}
	
	// TODO: 실제 구현에서는 새로운 카드 인스턴스 생성이 필요
	// 지금은 임시로 카드 ID에 "_upgraded" 접미사 추가
	upgradedCardID := fmt.Sprintf("%s_upgraded", originalCard.ID)
	
	// 덱에서 카드 교체
	playerState.Deck[cardIndex] = upgradedCardID
	
	return nil
}

// GetUpgradeCost 업그레이드 비용 계산
func (s *CardUpgradeServiceImpl) GetUpgradeCost(cardID string) int {
	// 카드 등급에 따른 업그레이드 비용
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return 100 // 기본 비용
	}
	
	switch card.Rarity {
	case domain.CardRarityCommon:
		return 50
	case domain.CardRarityRare:
		return 100
	case domain.CardRarityEpic:
		return 200
	case domain.CardRarityLegendary:
		return 500
	default:
		return 100
	}
}

// CanUpgradeCard 업그레이드 가능 여부
func (s *CardUpgradeServiceImpl) CanUpgradeCard(cardInstanceID string, playerState *domain.PlayerState) (bool, string) {
	// 카드가 덱에 있는지 확인
	hasCard := false
	for _, deckCardID := range playerState.Deck {
		if deckCardID == cardInstanceID {
			hasCard = true
			break
		}
	}
	
	if !hasCard {
		return false, "카드가 덱에 없습니다"
	}
	
	// 카드 정보 가져오기
	card, err := s.cardRepo.GetByID(cardInstanceID)
	if err != nil {
		return false, "카드 정보를 찾을 수 없습니다"
	}
	
	// 업그레이드 가능한 카드 타입인지 확인
	if !s.canUpgradeCardType(card) {
		return false, "이 카드는 업그레이드할 수 없습니다"
	}
	
	// 이미 업그레이드된 카드인지 확인
	if s.isAlreadyUpgraded(card) {
		return false, "이미 업그레이드된 카드입니다"
	}
	
	return true, ""
}

// canUpgradeCardType 업그레이드 가능한 카드 타입인지 확인
func (s *CardUpgradeServiceImpl) canUpgradeCardType(card *domain.Card) bool {
	// 기본적으로 모든 타입의 카드는 업그레이드 가능
	switch card.Type {
	case domain.CardTypeAction, domain.CardTypePower:
		return true
	default:
		return false
	}
}

// isAlreadyUpgraded 이미 업그레이드된 카드인지 확인
func (s *CardUpgradeServiceImpl) isAlreadyUpgraded(card *domain.Card) bool {
	// 카드 이름이나 ID에 업그레이드 표시가 있는지 확인
	// 실제 구현에서는 카드 인스턴스의 upgrade_level 필드를 확인
	return false // 임시로 false 반환
}

// createUpgradedCard 업그레이드된 카드 생성
func (s *CardUpgradeServiceImpl) createUpgradedCard(originalCard *domain.Card) *domain.Card {
	upgradedCard := *originalCard // 복사
	
	// 업그레이드 효과 적용
	switch originalCard.Type {
	case domain.CardTypeAction:
		// 공격 카드는 데미지 증가
		upgradedCard.BaseDamage += 3
		upgradedCard.Name = fmt.Sprintf("%s+", originalCard.Name)
		
	case domain.CardTypePower:
		// 파워 카드는 비용 감소 또는 효과 증가
		if originalCard.Cost > 1 {
			upgradedCard.Cost--
		} else {
			// 비용이 이미 낮으면 효과 증가
			if originalCard.BaseBlock > 0 {
				upgradedCard.BaseBlock += 3
			}
			if originalCard.DrawAmount > 0 {
				upgradedCard.DrawAmount += 1
			}
			if originalCard.BaseDamage > 0 {
				upgradedCard.BaseDamage += 2
			}
		}
		upgradedCard.Name = fmt.Sprintf("%s+", originalCard.Name)
	}
	
	// 설명 업데이트
	upgradedCard.Description = fmt.Sprintf("[업그레이드됨] %s", originalCard.Description)
	
	return &upgradedCard
}

// UpgradePreview 업그레이드 미리보기
func (s *CardUpgradeServiceImpl) UpgradePreview(cardID string) (*domain.Card, error) {
	originalCard, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("카드 정보 조회 실패: %w", err)
	}
	
	upgradedCard := s.createUpgradedCard(originalCard)
	return upgradedCard, nil
}

// GetUpgradeStats 업그레이드 통계
func (s *CardUpgradeServiceImpl) GetUpgradeStats(playerState *domain.PlayerState) map[string]interface{} {
	stats := map[string]interface{}{
		"total_cards":       len(playerState.Deck),
		"upgradeable_cards": 0,
		"upgraded_cards":    0,
		"upgrade_cost":      0,
	}
	
	upgradeableCards, _ := s.GetUpgradeableCards(playerState)
	stats["upgradeable_cards"] = len(upgradeableCards)
	
	// 총 업그레이드 비용 계산
	totalCost := 0
	for _, cardID := range upgradeableCards {
		totalCost += s.GetUpgradeCost(cardID)
	}
	stats["upgrade_cost"] = totalCost
	
	return stats
}