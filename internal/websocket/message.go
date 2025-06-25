package websocket

import (
	"time"
	"github.com/yourusername/pixel-game/internal/domain"
)

// MessageType WebSocket 메시지 타입
type MessageType string

const (
	// 연결 관련
	MessageTypeConnection    MessageType = "CONNECTION"
	MessageTypePing          MessageType = "PING"
	MessageTypePong          MessageType = "PONG"
	MessageTypeError         MessageType = "ERROR"

	// 게임 세션 관련
	MessageTypeSessionJoin   MessageType = "SESSION_JOIN"
	MessageTypeSessionJoined MessageType = "SESSION_JOINED"
	MessageTypeSessionLeave  MessageType = "SESSION_LEAVE"
	MessageTypeSessionLeft   MessageType = "SESSION_LEFT"

	// 게임 상태 관련
	MessageTypeGameState     MessageType = "GAME_STATE"
	MessageTypeGameAction    MessageType = "GAME_ACTION"
	MessageTypeGameUpdate    MessageType = "GAME_UPDATE"
	MessageTypeTurnStart     MessageType = "TURN_START"
	MessageTypeTurnEnd       MessageType = "TURN_END"

	// 카드 관련
	MessageTypeCardPlayed    MessageType = "CARD_PLAYED"
	MessageTypeCardDraw      MessageType = "CARD_DRAW"
	MessageTypeCardDiscard   MessageType = "CARD_DISCARD"

	// 전투 관련
	MessageTypeDamageDealt   MessageType = "DAMAGE_DEALT"
	MessageTypeShieldGained  MessageType = "SHIELD_GAINED"
	MessageTypeBuffApplied   MessageType = "BUFF_APPLIED"
	MessageTypeDebuffApplied MessageType = "DEBUFF_APPLIED"

	// 보상 관련
	MessageTypeRewardEarned  MessageType = "REWARD_EARNED"
	MessageTypeRewardSelect  MessageType = "REWARD_SELECT"

	// 시스템 메시지
	MessageTypeNotification MessageType = "NOTIFICATION"
	MessageTypeBroadcast    MessageType = "BROADCAST"
)

// Message WebSocket 메시지 구조체
type Message struct {
	Type      MessageType     `json:"type"`
	Data      interface{}     `json:"data"`
	Timestamp int64          `json:"timestamp"`
	MessageID string         `json:"message_id,omitempty"`
}

// NewMessage 새로운 메시지 생성
func NewMessage(msgType MessageType, data interface{}) Message {
	return Message{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// GameStateData 게임 상태 메시지 데이터
type GameStateData struct {
	SessionID    string      `json:"session_id"`
	CurrentTurn  int         `json:"current_turn"`
	TurnPhase    string      `json:"turn_phase"`
	PlayerState  interface{} `json:"player_state"`
	EnemyState   interface{} `json:"enemy_state"`
	GameState    interface{} `json:"game_state"`
}

// GameActionData 게임 액션 메시지 데이터
type GameActionData struct {
	SessionID  string      `json:"session_id"`
	ActionType string      `json:"action_type"`
	CardID     *string     `json:"card_id,omitempty"`
	TargetID   *string     `json:"target_id,omitempty"`
	ActionData interface{} `json:"action_data,omitempty"`
}

// CardPlayedData 카드 사용 메시지 데이터
type CardPlayedData struct {
	SessionID        string      `json:"session_id"`
	CardID           string      `json:"card_id"`
	PlayerID         int         `json:"player_id"`
	TargetID         *string     `json:"target_id,omitempty"`
	DamageDealt      int         `json:"damage_dealt,omitempty"`
	HealingDone      int         `json:"healing_done,omitempty"`
	ShieldGained     int         `json:"shield_gained,omitempty"`
	CardsDrawn       []string    `json:"cards_drawn,omitempty"`
	BuffsApplied     []domain.BuffState `json:"buffs_applied,omitempty"`
	DebuffsApplied   []domain.DebuffState `json:"debuffs_applied,omitempty"`
	Success          bool        `json:"success"`
	Messages         []string    `json:"messages,omitempty"`
	RemainingHand    []string    `json:"remaining_hand"`
}

// Effect 카드 효과 데이터
type Effect struct {
	Type        string      `json:"type"`
	Value       int         `json:"value"`
	Target      string      `json:"target"`
	Description string      `json:"description"`
	Success     bool        `json:"success"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// DamageData 데미지 메시지 데이터
type DamageData struct {
	SessionID     string `json:"session_id"`
	SourceID      string `json:"source_id"`
	TargetID      string `json:"target_id"`
	Damage        int    `json:"damage"`
	ActualDamage  int    `json:"actual_damage"`
	ShieldBlocked int    `json:"shield_blocked"`
	IsCritical    bool   `json:"is_critical"`
}

// BuffData 버프/디버프 메시지 데이터
type BuffData struct {
	SessionID   string `json:"session_id"`
	TargetID    string `json:"target_id"`
	BuffID      string `json:"buff_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Duration    int    `json:"duration"`
	IsDebuff    bool   `json:"is_debuff"`
}

// RewardData 보상 메시지 데이터
type RewardData struct {
	SessionID     string      `json:"session_id"`
	RewardBundle  interface{} `json:"reward_bundle"`
	FloorNumber   int         `json:"floor_number"`
	EnemyType     string      `json:"enemy_type"`
	HasChoices    bool        `json:"has_choices"`
}

// TurnData 턴 메시지 데이터
type TurnData struct {
	SessionID     string `json:"session_id"`
	TurnNumber    int    `json:"turn_number"`
	TurnPhase     string `json:"turn_phase"`
	CurrentPlayer string `json:"current_player"`
	TimeLimit     int    `json:"time_limit,omitempty"`
}

// NotificationData 알림 메시지 데이터
type NotificationData struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Type     string `json:"type"` // info, warning, error, success
	Duration int    `json:"duration,omitempty"` // 표시 시간 (초)
}

// ErrorData 에러 메시지 데이터
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// BroadcastData 브로드캐스트 메시지 데이터
type BroadcastData struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}