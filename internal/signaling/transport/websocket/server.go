package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене настройте CORS правильно
	},
}

type Client struct {
	ID   string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()

			h.clients[client] = true

			h.mu.Unlock()
			log.Printf("Client %s registered", client.ID)

			// Уведомить других о новом клиенте
			nitification, _ := json.Marshal(map[string]string{
				"type": "user_joined",
				"id":   client.ID,
			})
			h.broadcastMessage(nitification, client)

		case client := <-h.unregister:
			h.mu.Lock()

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.mu.Unlock()

				log.Printf("Client %s unregistered", client.ID)

				// Уведомить других о отключении клиента
				nitification, _ := json.Marshal(map[string]string{
					"type": "user_left",
					"id":   client.ID,
				})
				h.broadcastMessage(nitification, client)
			} else {
				h.mu.Unlock()
			}
		case message := <-h.broadcast:
			h.broadcastMessage(message, nil)

		}
	}
}

func (h *Hub) broadcastMessage(message []byte, exclude *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client != exclude {
			select {
			case client.send <- message:
			default:
				// клиентский канал отправки заполнен, закроем его
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка чтения: %v", err)
			}
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		switch msg["type"] {
		case "broadcast":
			c.hub.broadcast <- message
		case "ping":
			pong, _ := json.Marshal(map[string]string{"type": "pong"})
			c.send <- pong
		default:
			c.send <- message
		}

	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
