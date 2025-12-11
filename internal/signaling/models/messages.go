package models

import (
	"encoding/json"

	"github.com/stazoloto/sfu-mediaserver/internal/sfu/domain/models"
	"github.com/stazoloto/sfu-mediaserver/internal/sfu/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/sfu/domain/vo/room"
)

const (
	TypeJoin      = "join"
	TypeLeave     = "leave"
	TypeOffer     = "offer"
	TypeAnswer    = "answer"
	TypeCandidate = "candidate"
	TypeError     = "error"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	Room models.Room     `json:"room,omitempty"`
}

// Для присоединения к комнате
type JoinMessage struct {
	RoomID  room.ID `json:"room_id"`
	PeerID  peer.ID `json:"peer_id"`
	IsOffer bool    `json:"is_offer"` // является ли создателем комнаты
}

// WebRTC SDP сообщение
type SDPMessage struct {
	SDP string `json:"sdp"`
	To  string `json:"to,omitempty"`
}

// ICE кандидат
type ICECandidate struct {
	Candidate     string `json:"candidate"`
	SDPMLineIndex int    `json:"sdpMLineIndex"` // Индекс медиа-потока
	SDPMid        string `json:"sdpMid"`        // Идентификатор медиа-потока
	To            string `json:"to,omitempty"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
