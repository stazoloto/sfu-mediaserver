package models

import (
	"time"

	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
)

type Peer struct {
	ID        peer.ID   `json:"id"`
	Name      string    `json:"name"`
	Role      peer.Role `json:"role"`
	SessionID string    `json:"session_id"`
	JoinTime  time.Time `json:"join_time"`
	Metadata  []byte    `json:"metadata"`
}

func NewPeer(
	name string,
	role peer.Role,
	sessionID string,
	joinTime time.Time,
	metadata []byte,
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
		Metadata:  metadata,
	}
}
