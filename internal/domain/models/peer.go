package models

import (
	"time"

	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
)

// Peer - представление пользователя в конкретной комнате
type Peer struct {
	ID        peer.ID   `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Role      peer.Role `json:"role" db:"role"` // "member", "owner"
	SessionID string    `json:"session_id" db:"session_id"`
	JoinTime  time.Time `json:"join_time" db:"join_time"`
}

func NewPeer(
	name string,
	role peer.Role,
	sessionID string,
	joinTime time.Time,
) *Peer {
	id, err := peer.NewPeerID()
	if err != nil {
		panic(err)
	}
	return &Peer{
		ID:        id,
		Name:      name,
		Role:      role,
		SessionID: sessionID,
		JoinTime:  joinTime,
	}
}
