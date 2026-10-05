package hub

import (
	"encoding/json"
	"sync"

	"ThinkTalk/application/chat/api/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Conn struct {
	UserId int64
	Ws     *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	mu        sync.RWMutex
	conns     map[int64]*Conn
	OnMessage func(userId int64, msg *types.WsInMessage)
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[int64]*Conn),
	}
}

func (h *Hub) Register(userId int64, ws *websocket.Conn) *Conn {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, ok := h.conns[userId]; ok {
		close(old.Send)
		old.Ws.Close()
	}

	conn := &Conn{
		UserId: userId,
		Ws:     ws,
		Send:   make(chan []byte, 64),
	}
	h.conns[userId] = conn
	return conn
}

func (h *Hub) Unregister(userId int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.conns[userId]; ok {
		close(conn.Send)
		conn.Ws.Close()
		delete(h.conns, userId)
	}
}

func (h *Hub) SendToUser(userId int64, msg *types.WsOutMessage) {
	h.mu.RLock()
	conn, ok := h.conns[userId]
	h.mu.RUnlock()

	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logx.Errorf("[Hub] marshal msg err: %v", err)
		return
	}

	select {
	case conn.Send <- data:
	default:
	}
}

func (h *Hub) ReadPump(conn *Conn) {
	defer func() {
		h.Unregister(conn.UserId)
	}()

	for {
		_, message, err := conn.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logx.Errorf("[Hub] ws read error: %v userId: %d", err, conn.UserId)
			}
			break
		}

		var msg types.WsInMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			logx.Errorf("[Hub] unmarshal ws msg err: %v userId: %d", err, conn.UserId)
			continue
		}

		if h.OnMessage != nil {
			h.OnMessage(conn.UserId, &msg)
		}
	}
}

func (h *Hub) WritePump(conn *Conn) {
	defer func() {
		conn.Ws.Close()
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			if !ok {
				conn.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.Ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}
