package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	upgrader websocket.Upgrader

	mu      sync.RWMutex
	clients map[string]*WSCLient // clientID -> conn

	controller   Controller
	onDisconnect func(clientID string)
}

type WSCLient struct {
	id   string
	conn *websocket.Conn
	send chan []byte
}

func (c *WSCLient) writeLoop() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// NewHub создаёт WebSocket hub
func NewHub() *Hub {
	return &Hub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[string]*WSCLient),
	}
}

// SetController подключает inbound adapter (WSController)
func (h *Hub) SetController(c Controller) {
	h.controller = c
}

// ServeHTTP — entrypoint для net/http
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// client_id передаётся как query param
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(4001, "client_id is required"),
		)
		return
	}

	h.register(clientID, conn)
	defer h.unregister(clientID)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("ws closed:", clientID)
			} else {
				log.Println("ws read error:", err)
			}

			if h.onDisconnect != nil {
				h.onDisconnect(clientID)
			}
			return
		}

		if h.controller != nil {
			if err := h.controller.Handle(data); err != nil {
				log.Println("controller error:", err)
			}
		}
	}
}

func (h *Hub) register(clientID string, conn *websocket.Conn) {
	client := &WSCLient{
		id:   clientID,
		conn: conn,
		send: make(chan []byte, 256),
	}
	h.mu.Lock()
	h.clients[clientID] = client
	h.mu.Unlock()

	go client.writeLoop()

	log.Println("WS CONNECT:", clientID)
}

func (h *Hub) unregister(clientID string) {
	h.mu.Lock()
	client := h.clients[clientID]
	delete(h.clients, clientID)
	h.mu.Unlock()

	if client != nil {
		close(client.send)
	}

	log.Println("WS DISCONNECT:", clientID)
}

// Send отправляет raw JSON конкретному клиенту
func (h *Hub) Send(clientID string, data []byte) error {
	h.mu.RLock()
	client := h.clients[clientID]
	h.mu.RUnlock()

	if client == nil {
		return nil
	}

	select {
	case client.send <- data:
	default:
	}
	return nil
}

func (h *Hub) SetOnDisconnect(fn func(string)) {
	h.onDisconnect = fn
}
