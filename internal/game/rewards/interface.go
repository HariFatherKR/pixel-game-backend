package rewards

import (
	"github.com/yourusername/pixel-game/internal/domain"
)

// RewardType 보상 타입 상수
type RewardType string

const (
	RewardTypeGold     RewardType = "GOLD"     // 골드 보상
	RewardTypeCard     RewardType = "CARD"     // 카드 보상
	RewardTypeRelic    RewardType = "RELIC"    // 유물 보상
	RewardTypePotion   RewardType = "POTION"   // 포션 보상
	RewardTypeUpgrade  RewardType = "UPGRADE"  // 카드 업그레이드
	RewardTypeRemove   RewardType = "REMOVE"   // 카드 제거
	RewardTypeHealth   RewardType = "HEALTH"   // 체력 회복
)

// RewardRarity 보상 등급
type RewardRarity string

const (
	RewardRarityCommon    RewardRarity = "COMMON"
	RewardRarityRare      RewardRarity = "RARE"
	RewardRarityEpic      RewardRarity = "EPIC"
	RewardRarityLegendary RewardRarity = "LEGENDARY"
)

// Reward 개별 보상 정보
type Reward struct {
	ID          string                 `json:"id"`
	Type        RewardType             `json:"type"`
	Rarity      RewardRarity           `json:"rarity"`
	Value       int                    `json:"value"`       // 골드량, 체력량 등
	ItemID      string                 `json:"item_id"`     // 카드 ID, 유물 ID 등
	Name        string                 `json:"name"`        // 보상 이름
	Description string                 `json:"description"` // 보상 설명
	ImageURL    string                 `json:"image_url"`   // 이미지 URL
	Metadata    map[string]interface{} `json:"metadata"`    // 추가 메타데이터
}

// RewardBundle 보상 묶음 (전투 승리 시 받는 전체 보상)
type RewardBundle struct {
	ID           string    `json:"id"`
	SourceType   string    `json:"source_type"`   // "COMBAT", "EVENT", "BOSS"
	SourceID     string    `json:"source_id"`     // 적 ID, 이벤트 ID 등
	FloorNumber  int       `json:"floor_number"`  // 층 번호
	BaseRewards  []Reward  `json:"base_rewards"`  // 기본 보상 (항상 주어짐)
	ChoiceRewards []Reward `json:"choice_rewards"` // 선택 보상 (하나만 선택)
	IsCompleted  bool      `json:"is_completed"`  // 보상 수령 완료 여부
}

// RewardContext 보상 생성을 위한 컨텍스트
type RewardContext struct {
	FloorNumber   int                    `json:"floor_number"`
	EnemyType     string                 `json:"enemy_type"`
	PlayerLevel   int                    `json:"player_level"`
	GameMode      domain.GameMode        `json:"game_mode"`
	PlayerState   *domain.PlayerState    `json:"player_state"`
	GameState     *domain.GameState      `json:"game_state"`
	DifficultyMod float64                `json:"difficulty_mod"` // 난이도 배율
	BonusFactors  map[string]interface{} `json:"bonus_factors"`  // 추가 보너스 요소들
}

// RewardGenerator 보상 생성 인터페이스
type RewardGenerator interface {
	// GenerateRewards 보상 묶음 생성
	GenerateRewards(ctx *RewardContext) (*RewardBundle, error)
	
	// GenerateCardRewards 카드 보상 생성
	GenerateCardRewards(ctx *RewardContext, count int) ([]Reward, error)
	
	// GenerateGoldReward 골드 보상 생성
	GenerateGoldReward(ctx *RewardContext) (*Reward, error)
	
	// GenerateRelicReward 유물 보상 생성 (보스전 등)
	GenerateRelicReward(ctx *RewardContext) (*Reward, error)
	
	// CalculateRewardValue 보상 가치 계산
	CalculateRewardValue(rewardType RewardType, ctx *RewardContext) int
}

// RewardManager 보상 관리 시스템
type RewardManager interface {
	// ProcessRewards 보상 처리
	ProcessRewards(
		sessionID string,
		playerState *domain.PlayerState,
		gameState *domain.GameState,
		ctx *RewardContext,
	) (*RewardBundle, error)
	
	// ApplyReward 개별 보상 적용
	ApplyReward(
		sessionID string,
		playerState *domain.PlayerState,
		gameState *domain.GameState,
		reward *Reward,
	) error
	
	// ValidateRewardChoice 보상 선택 유효성 검사
	ValidateRewardChoice(bundleID string, rewardID string) (bool, string)
	
	// CompleteRewardSelection 보상 선택 완료
	CompleteRewardSelection(
		sessionID string,
		bundleID string,
		selectedRewardIDs []string,
		playerState *domain.PlayerState,
		gameState *domain.GameState,
	) error
	
	// GetPendingRewards 대기 중인 보상 목록
	GetPendingRewards(sessionID string) ([]*RewardBundle, error)
	
	// GetRewardHistory 보상 히스토리
	GetRewardHistory(sessionID string) ([]*RewardBundle, error)
	
	// CalculateSessionRewards 세션 보상 통계
	CalculateSessionRewards(sessionID string) (map[string]interface{}, error)
}

// RewardRepository 보상 데이터 저장소 인터페이스
type RewardRepository interface {
	// SaveRewardBundle 보상 묶음 저장
	SaveRewardBundle(sessionID string, bundle *RewardBundle) error
	
	// GetRewardBundle 보상 묶음 조회
	GetRewardBundle(sessionID string, bundleID string) (*RewardBundle, error)
	
	// GetPendingRewards 대기 중인 보상 목록
	GetPendingRewards(sessionID string) ([]*RewardBundle, error)
	
	// MarkRewardCompleted 보상 완료 처리
	MarkRewardCompleted(sessionID string, bundleID string) error
	
	// GetRewardHistory 보상 히스토리
	GetRewardHistory(sessionID string) ([]*RewardBundle, error)
}

// CardUpgradeService 카드 업그레이드 서비스
type CardUpgradeService interface {
	// GetUpgradeableCards 업그레이드 가능한 카드 목록
	GetUpgradeableCards(playerState *domain.PlayerState) ([]string, error)
	
	// UpgradeCard 카드 업그레이드
	UpgradeCard(cardInstanceID string, playerState *domain.PlayerState) error
	
	// GetUpgradeCost 업그레이드 비용 계산
	GetUpgradeCost(cardID string) int
	
	// CanUpgradeCard 업그레이드 가능 여부
	CanUpgradeCard(cardInstanceID string, playerState *domain.PlayerState) (bool, string)
	
	// UpgradePreview 업그레이드 미리보기
	UpgradePreview(cardID string) (*domain.Card, error)
	
	// GetUpgradeStats 업그레이드 통계
	GetUpgradeStats(playerState *domain.PlayerState) map[string]interface{}
}

// RewardEvent 보상 이벤트
type RewardEvent struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id"`
	Reward    *Reward                `json:"reward"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp int64                  `json:"timestamp"`
}

// RewardEventType 보상 이벤트 타입
const (
	RewardEventTypeGenerated = "REWARD_GENERATED" // 보상 생성됨
	RewardEventTypeSelected  = "REWARD_SELECTED"  // 보상 선택됨
	RewardEventTypeApplied   = "REWARD_APPLIED"   // 보상 적용됨
	RewardEventTypeSkipped   = "REWARD_SKIPPED"   // 보상 건너뜀
)