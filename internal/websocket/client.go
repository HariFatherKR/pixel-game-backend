package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 클라이언트에게 메시지를 보낼 때의 시간 제한
	writeWait = 10 * time.Second

	// 클라이언트로부터 다음 pong 메시지를 기다리는 시간
	pongWait = 60 * time.Second

	// 이 기간 내에 pong 메시지를 전송해야 함 (pongWait보다 작아야 함)
	pingPeriod = (pongWait * 9) / 10

	// 클라이언트로부터 메시지를 읽을 때의 최대 크기
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 개발 환경에서는 모든 origin 허용, 프로덕션에서는 검증 필요
		return true
	},
}

// Client WebSocket 클라이언트
type Client struct {
	// WebSocket 연결
	conn *websocket.Conn

	// 메시지를 허브로 전송하는 채널
	hub *Hub

	// 클라이언트에게 전송할 메시지들의 버퍼 채널
	send chan []byte

	// 사용자 ID
	UserID int

	// 현재 게임 세션 ID
	SessionID string

	// 마지막 활동 시간
	LastActivity time.Time
}

// NewClient 새로운 클라이언트 생성
func NewClient(hub *Hub, conn *websocket.Conn, userID int, sessionID string) *Client {
	return &Client{
		conn:         conn,
		hub:          hub,
		send:         make(chan []byte, 256),
		UserID:       userID,
		SessionID:    sessionID,
		LastActivity: time.Now(),
	}
}

// readPump 클라이언트로부터 WebSocket 연결의 메시지를 읽어들임
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.LastActivity = time.Now()
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 오류: %v", err)
			}
			break
		}

		c.LastActivity = time.Now()

		// 받은 메시지 처리
		var message Message
		if err := json.Unmarshal(messageData, &message); err != nil {
			log.Printf("메시지 파싱 오류: %v", err)
			continue
		}

		// 메시지 타입별 처리
		c.handleMessage(&message)
	}
}

// writePump 허브로부터 메시지를 받아 WebSocket 연결로 전송
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 허브가 채널을 닫음
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 대기 중인 채팅 메시지들을 더 추가
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 받은 메시지 처리
func (c *Client) handleMessage(message *Message) {
	switch message.Type {
	case MessageTypePing:
		c.handlePing(message)
	case MessageTypeGameAction:
		c.handleGameAction(message)
	case MessageTypeSessionJoin:
		c.handleSessionJoin(message)
	case MessageTypeSessionLeave:
		c.handleSessionLeave(message)
	default:
		log.Printf("알 수 없는 메시지 타입: %s", message.Type)
	}
}

// handlePing Ping 메시지 처리
func (c *Client) handlePing(message *Message) {
	response := Message{
		Type: MessageTypePong,
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"server":    "pixel-game-backend",
		},
	}
	c.SendMessage(response)
}

// handleGameAction 게임 액션 메시지 처리
func (c *Client) handleGameAction(message *Message) {
	// 게임 액션은 HTTP API를 통해 처리되므로 여기서는 확인만
	log.Printf("게임 액션 수신 - UserID: %d, SessionID: %s", c.UserID, c.SessionID)
	
	// 실제 게임 로직은 HTTP API 핸들러에서 처리되고
	// 결과는 WebSocket을 통해 브로드캐스트됨
}

// handleSessionJoin 게임 세션 참가 처리
func (c *Client) handleSessionJoin(message *Message) {
	if data, ok := message.Data.(map[string]interface{}); ok {
		if sessionID, ok := data["session_id"].(string); ok {
			c.SessionID = sessionID
			log.Printf("클라이언트가 세션에 참가함 - UserID: %d, SessionID: %s", c.UserID, sessionID)
			
			// 세션 참가 성공 응답
			response := Message{
				Type: MessageTypeSessionJoined,
				Data: map[string]interface{}{
					"session_id": sessionID,
					"status":     "joined",
					"message":    "게임 세션에 참가했습니다",
				},
			}
			c.SendMessage(response)
		}
	}
}

// handleSessionLeave 게임 세션 떠나기 처리
func (c *Client) handleSessionLeave(message *Message) {
	oldSessionID := c.SessionID
	c.SessionID = ""
	
	log.Printf("클라이언트가 세션을 떠남 - UserID: %d, SessionID: %s", c.UserID, oldSessionID)
	
	// 세션 떠나기 성공 응답
	response := Message{
		Type: MessageTypeSessionLeft,
		Data: map[string]interface{}{
			"session_id": oldSessionID,
			"status":     "left",
			"message":    "게임 세션에서 나갔습니다",
		},
	}
	c.SendMessage(response)
}

// SendMessage 클라이언트에게 메시지 전송
func (c *Client) SendMessage(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("메시지 직렬화 실패: %v", err)
		return
	}

	select {
	case c.send <- data:
	default:
		log.Printf("클라이언트 %d 메시지 전송 실패: 채널 가득함", c.UserID)
	}
}

// ServeWS WebSocket 연결을 처리하고 클라이언트를 허브에 등록
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, userID int, sessionID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 업그레이드 실패: %v", err)
		return
	}

	client := NewClient(hub, conn, userID, sessionID)
	client.hub.register <- client

	// 각각을 별도의 고루틴에서 실행
	go client.writePump()
	go client.readPump()
}