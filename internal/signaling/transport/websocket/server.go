package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[string]*Client
	rooms      map[string][]string
	register   chan *Client          // канал для регистрации новых клиентов
	unregister chan *Client          // канал для удаления клиентов
	broadcast  chan BroadcastMessage // канал для широковещательных сообщений
	mu         sync.RWMutex          // для безопасного доступа из горутин
}

type Client struct {
	conn   *websocket.Conn
	peerID string
	roomID string
	send   chan []byte
	hub    *Hub
}

type BroadcastMessage struct {
	RoomID  string
	Exclude string
	Message []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		rooms:      make(map[string][]string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {}
	}
}
