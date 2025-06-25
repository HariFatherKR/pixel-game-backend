package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub WebSocket 연결 관리 허브
type Hub struct {
	// 등록된 클라이언트들
	clients map[*Client]bool

	// 클라이언트 등록 채널
	register chan *Client

	// 클라이언트 해제 채널
	unregister chan *Client

	// 클라이언트들에게 메시지 브로드캐스트하는 채널
	broadcast chan []byte

	// 특정 사용자에게 메시지 전송하는 채널
	sendToUser chan *UserMessage

	// 특정 게임 세션에 메시지 전송하는 채널
	sendToSession chan *SessionMessage

	// 사용자 ID별 클라이언트 매핑
	userClients map[int]*Client

	// 게임 세션별 클라이언트 매핑
	sessionClients map[string][]*Client

	// 뮤텍스
	mu sync.RWMutex
}

// UserMessage 특정 사용자에게 보내는 메시지
type UserMessage struct {
	UserID  int    `json:"user_id"`
	Message []byte `json:"message"`
}

// SessionMessage 특정 게임 세션에 보내는 메시지
type SessionMessage struct {
	SessionID string `json:"session_id"`
	Message   []byte `json:"message"`
}

// NewHub 새로운 허브 생성
func NewHub() *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		sendToUser:     make(chan *UserMessage),
		sendToSession:  make(chan *SessionMessage),
		userClients:    make(map[int]*Client),
		sessionClients: make(map[string][]*Client),
	}
}

// Run 허브 실행
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case userMsg := <-h.sendToUser:
			h.sendMessageToUser(userMsg)

		case sessionMsg := <-h.sendToSession:
			h.sendMessageToSession(sessionMsg)
		}
	}
}

// registerClient 클라이언트 등록
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	
	// 사용자 ID별 매핑
	if client.UserID != 0 {
		h.userClients[client.UserID] = client
	}

	// 게임 세션별 매핑
	if client.SessionID != "" {
		h.sessionClients[client.SessionID] = append(h.sessionClients[client.SessionID], client)
	}

	log.Printf("클라이언트 연결됨 - UserID: %d, SessionID: %s", client.UserID, client.SessionID)

	// 연결 성공 메시지 전송
	welcomeMsg := Message{
		Type: MessageTypeConnection,
		Data: map[string]interface{}{
			"status":     "connected",
			"message":    "WebSocket 연결이 성공했습니다",
			"user_id":    client.UserID,
			"session_id": client.SessionID,
		},
	}
	
	if msgData, err := json.Marshal(welcomeMsg); err == nil {
		select {
		case client.send <- msgData:
		default:
			h.forceUnregisterClient(client)
		}
	}
}

// unregisterClient 클라이언트 해제
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		h.forceUnregisterClientLocked(client)
	}
}

// forceUnregisterClient 강제 클라이언트 해제
func (h *Hub) forceUnregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.forceUnregisterClientLocked(client)
}

// forceUnregisterClientLocked 락이 이미 걸린 상태에서 클라이언트 해제
func (h *Hub) forceUnregisterClientLocked(client *Client) {
	delete(h.clients, client)
	close(client.send)

	// 사용자 ID별 매핑에서 제거
	if client.UserID != 0 {
		delete(h.userClients, client.UserID)
	}

	// 게임 세션별 매핑에서 제거
	if client.SessionID != "" {
		if clients, ok := h.sessionClients[client.SessionID]; ok {
			for i, c := range clients {
				if c == client {
					h.sessionClients[client.SessionID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			// 세션에 클라이언트가 없으면 세션 삭제
			if len(h.sessionClients[client.SessionID]) == 0 {
				delete(h.sessionClients, client.SessionID)
			}
		}
	}

	log.Printf("클라이언트 연결 해제됨 - UserID: %d, SessionID: %s", client.UserID, client.SessionID)
}

// broadcastMessage 모든 클라이언트에게 메시지 브로드캐스트
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			go h.forceUnregisterClient(client)
		}
	}
}

// sendMessageToUser 특정 사용자에게 메시지 전송
func (h *Hub) sendMessageToUser(userMsg *UserMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.userClients[userMsg.UserID]; ok {
		select {
		case client.send <- userMsg.Message:
		default:
			go h.forceUnregisterClient(client)
		}
	}
}

// sendMessageToSession 특정 게임 세션의 모든 클라이언트에게 메시지 전송
func (h *Hub) sendMessageToSession(sessionMsg *SessionMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.sessionClients[sessionMsg.SessionID]; ok {
		for _, client := range clients {
			select {
			case client.send <- sessionMsg.Message:
			default:
				go h.forceUnregisterClient(client)
			}
		}
	}
}

// SendToUser 외부에서 특정 사용자에게 메시지 전송
func (h *Hub) SendToUser(userID int, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("메시지 직렬화 실패: %v", err)
		return
	}

	select {
	case h.sendToUser <- &UserMessage{UserID: userID, Message: data}:
	default:
		log.Printf("사용자 %d에게 메시지 전송 실패: 채널 가득참", userID)
	}
}

// SendToSession 외부에서 특정 게임 세션에 메시지 전송
func (h *Hub) SendToSession(sessionID string, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("메시지 직렬화 실패: %v", err)
		return
	}

	select {
	case h.sendToSession <- &SessionMessage{SessionID: sessionID, Message: data}:
	default:
		log.Printf("세션 %s에 메시지 전송 실패: 채널 가득참", sessionID)
	}
}

// Broadcast 외부에서 모든 클라이언트에게 브로드캐스트
func (h *Hub) Broadcast(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("메시지 직렬화 실패: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		log.Printf("브로드캐스트 실패: 채널 가득참")
	}
}

// GetConnectedUsers 연결된 사용자 수 반환
func (h *Hub) GetConnectedUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients)
}

// GetActiveSessions 활성 게임 세션 수 반환
func (h *Hub) GetActiveSessions() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.sessionClients)
}

// IsUserConnected 특정 사용자가 연결되어 있는지 확인
func (h *Hub) IsUserConnected(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userClients[userID]
	return ok
}

// IsSessionActive 특정 게임 세션이 활성 상태인지 확인
func (h *Hub) IsSessionActive(sessionID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.sessionClients[sessionID]
	return ok && len(clients) > 0
}