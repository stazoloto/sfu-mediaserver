package models

import (
	"time"

	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/room"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/session"
)

// Session - одно подключение участника.
// У одного Peer может быть несколько сессий (например, разные устройства)
type Session struct {
	ID     session.ID `json:"id"`
	PeerID peer.ID    `json:"peer_id"`
	RoomID room.ID    `json:"room_id"`

	ICECandidates   []session.ICECandidate  `json:"ice_candidates"`
	SignalingState  session.SignalingState  `json:"signaling_state"`
	ConnectionState session.ConnectionState `json:"connection_state"`

	// WebRTC состояние
	PeerConnection interface{}   `json:"-"`
	CreatedAt      time.Time     `json:"created_at"`
	LastActivity   time.Time     `json:"last_activity"`
	Duration       time.Duration `json:"duration"`
}

func NewSession(
	peerID peer.ID,
	roomID room.ID,
	iCECandidates []session.ICECandidate,
	signalingState session.SignalingState,
	connectionState session.ConnectionState,
	peerConnection interface{},
	createdAt time.Time,
	lastActivity time.Time,
	duration time.Duration,
) *Session {
	id, err := session.NewSessionID()
	if err != nil {
		panic(err)
	}
	return &Session{
		ID:              id,
		PeerID:          peerID,
		RoomID:          roomID,
		ICECandidates:   iCECandidates,
		SignalingState:  signalingState,
		ConnectionState: connectionState,
		PeerConnection:  peerConnection,
		CreatedAt:       createdAt,
		LastActivity:    lastActivity,
		Duration:        duration,
	}
}
