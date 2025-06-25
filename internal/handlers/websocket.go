package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/auth"
	"github.com/yourusername/pixel-game/internal/websocket"
)

// WebSocketHandler WebSocket 핸들러
type WebSocketHandler struct {
	hub        *websocket.Hub
	jwtManager *auth.JWTManager
}

// NewWebSocketHandler WebSocket 핸들러 생성
func NewWebSocketHandler(hub *websocket.Hub, jwtManager *auth.JWTManager) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		jwtManager: jwtManager,
	}
}

// HandleWebSocket WebSocket 연결 처리
// @Summary WebSocket 연결
// @Description 실시간 게임 통신을 위한 WebSocket 연결을 설정합니다
// @Tags WebSocket
// @Param Authorization header string true "Bearer 토큰"
// @Param session_id query string false "게임 세션 ID"
// @Success 101 "WebSocket 연결 성공"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 토큰에서 직접 사용자 ID 추출 (WebSocket은 헤더에서 토큰을 가져와야 함)
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "토큰이 필요합니다"})
		return
	}

	// "Bearer " 접두사 제거
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// 토큰 검증
	claims, err := h.jwtManager.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰"})
		return
	}

	userID := claims.UserID

	// 선택적 세션 ID
	sessionID := c.Query("session_id")

	// WebSocket 연결 업그레이드 및 클라이언트 등록
	websocket.ServeWS(h.hub, c.Writer, c.Request, userID, sessionID)
}

// GetWebSocketStats WebSocket 연결 통계 조회
// @Summary WebSocket 연결 통계
// @Description 현재 WebSocket 연결 상태와 통계를 조회합니다
// @Tags WebSocket
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /ws/stats [get]
func (h *WebSocketHandler) GetWebSocketStats(c *gin.Context) {
	stats := map[string]interface{}{
		"connected_users":   h.hub.GetConnectedUsers(),
		"active_sessions":   h.hub.GetActiveSessions(),
		"server_time":       gin.H{"timestamp": "now"},
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
		"status": "online",
	})
}

// SendNotification 특정 사용자에게 알림 전송
// @Summary 사용자 알림 전송
// @Description 특정 사용자에게 WebSocket을 통해 알림을 전송합니다
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param user_id path int true "사용자 ID"
// @Param notification body NotificationRequest true "알림 내용"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /ws/users/{user_id}/notify [post]
func (h *WebSocketHandler) SendNotification(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 사용자 ID"})
		return
	}

	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식"})
		return
	}

	// 사용자가 연결되어 있는지 확인
	if !h.hub.IsUserConnected(userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자가 연결되어 있지 않습니다"})
		return
	}

	// 알림 메시지 생성
	notification := websocket.NewMessage(websocket.MessageTypeNotification, websocket.NotificationData{
		Title:    req.Title,
		Message:  req.Message,
		Type:     req.Type,
		Duration: req.Duration,
	})

	// 사용자에게 알림 전송
	h.hub.SendToUser(userID, notification)

	c.JSON(http.StatusOK, gin.H{
		"message": "알림이 전송되었습니다",
		"user_id": userID,
	})
}

// SendSessionMessage 게임 세션에 메시지 전송
// @Summary 게임 세션 메시지 전송
// @Description 특정 게임 세션의 모든 참가자에게 메시지를 전송합니다
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param session_id path string true "게임 세션 ID"
// @Param message body SessionMessageRequest true "메시지 내용"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /ws/sessions/{session_id}/message [post]
func (h *WebSocketHandler) SendSessionMessage(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req SessionMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식"})
		return
	}

	// 세션이 활성 상태인지 확인
	if !h.hub.IsSessionActive(sessionID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "세션이 활성 상태가 아닙니다"})
		return
	}

	// 메시지 생성
	message := websocket.NewMessage(websocket.MessageType(req.Type), req.Data)

	// 세션에 메시지 전송
	h.hub.SendToSession(sessionID, message)

	c.JSON(http.StatusOK, gin.H{
		"message":    "세션에 메시지가 전송되었습니다",
		"session_id": sessionID,
	})
}

// SendBroadcast 전체 브로드캐스트 메시지 전송
// @Summary 전체 브로드캐스트
// @Description 연결된 모든 클라이언트에게 브로드캐스트 메시지를 전송합니다
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param broadcast body BroadcastRequest true "브로드캐스트 내용"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /ws/broadcast [post]
func (h *WebSocketHandler) SendBroadcast(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식"})
		return
	}

	// 브로드캐스트 메시지 생성
	broadcast := websocket.NewMessage(websocket.MessageTypeBroadcast, websocket.BroadcastData{
		Title:     req.Title,
		Message:   req.Message,
		Type:      req.Type,
		Timestamp: websocket.NewMessage(websocket.MessageTypeBroadcast, nil).Timestamp,
	})

	// 전체 브로드캐스트
	h.hub.Broadcast(broadcast)

	c.JSON(http.StatusOK, gin.H{
		"message": "브로드캐스트가 전송되었습니다",
		"recipients": h.hub.GetConnectedUsers(),
	})
}

// NotificationRequest 알림 요청 구조체
type NotificationRequest struct {
	Title    string `json:"title" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Type     string `json:"type" binding:"required"` // info, warning, error, success
	Duration int    `json:"duration"`
}

// SessionMessageRequest 세션 메시지 요청 구조체
type SessionMessageRequest struct {
	Type string      `json:"type" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

// BroadcastRequest 브로드캐스트 요청 구조체
type BroadcastRequest struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

// RegisterRoutes WebSocket 관련 라우트 등록
func (h *WebSocketHandler) RegisterRoutes(router *gin.RouterGroup) {
	// WebSocket 연결
	router.GET("/ws", h.HandleWebSocket)
	
	// WebSocket 관리 API (인증 필요)
	wsGroup := router.Group("/ws")
	wsGroup.GET("/stats", h.GetWebSocketStats)
	wsGroup.POST("/users/:user_id/notify", h.SendNotification)
	wsGroup.POST("/sessions/:session_id/message", h.SendSessionMessage)
	wsGroup.POST("/broadcast", h.SendBroadcast)
}