package models

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/stazoloto/sfu-mediaserver/internal/sfu/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/sfu/domain/vo/room"
)

// Peer - представление пользователя в конкретной комнате
type Peer struct {
	ID       peer.ID          `json:"id"`
	RoomID   room.ID          `json:"room_id"`
	Name     string           `json:"name"`
	Role     peer.Role        `json:"role"` // "member", "owner"
	Conn     *WebsocketClient `json:"-"`
	JoinTime time.Time        `json:"join_time"`
}

func NewPeer(
	roomID room.ID,
	name string,
	role peer.Role,
) *Peer {
	id, err := peer.NewPeerID()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	return &Peer{
		ID:       id,
		RoomID:   roomID,
		Name:     name,
		Role:     role,
		JoinTime: now,
	}
}
