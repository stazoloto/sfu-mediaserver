package types

import (
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/room"
)

const (
	TypeJoin      = "join"
	TypeLeave     = "leave"
	TypeOffer     = "offer"
	TypeAnswer    = "answer"
	TypeCandidate = "candidate"
	TypeError     = "error"
)

type SignalMessage struct {
	Type    string  `json:"type"`
	RoomID  room.ID `json:"room_id,omitempty"`
	PeerID  peer.ID `json:"peer_id,omitempty"`
	Payload any     `json:"payload,omitempty"`
}

// Для присоединения к комнате
type JoinRequest struct {
	RoomID room.ID `json:"room_id"`
	PeerID peer.ID `json:"peer_id"`
}

// WebRTC оффер
type OfferPayload struct {
	SDP  string `json:"sdp"`
	From string `json:"from"`
	To   string `json:"to"`
}

// WebRTC ответ
type AnswerPayload struct {
	SDP  string `json:"sdp"`
	From string `json:"from"`
	To   string `json:"to"`
}

// ICE кандидат
type ICECandidate struct {
	Candidate     string `json:"candidate"`
	SDPMLineIndex int    `json:"sdpMLineIndex"` // Индекс медиа-потока
	SDPMid        string `json:"sdpMid"`        // Идентификатор медиа-потока
	From          string `json:"from"`
	To            string `json:"to,omitempty"`
}
