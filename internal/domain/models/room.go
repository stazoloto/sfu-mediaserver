package models

import (
	"time"

	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/room"
)

// Room - комната
type Room struct {
	ID        room.ID           `json:"id"`
	Name      *string           `json:"name,omitempty"`
	OwnerID   peer.ID           `json:"owner_id"`
	Peers     map[string]*Peer  `json:"peers"`
	Tracks    map[string]*Track `json:"tracks"`
	CreatedAt time.Time         `json:"created_at"`
	DeletedAt time.Time         `json:"deleted_at"`
	IsActive  bool              `json:"is_active"`
}

func NewRoom(
	peers map[string]*Peer,
	tracks map[string]*Track,
) (*Room, error) {
	id, err := room.NewRoomID()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return &Room{
		ID:        id,
		Peers:     peers,
		Tracks:    tracks,
		CreatedAt: now,
		DeletedAt: now.AddDate(0, 1, 0),
		IsActive:  true,
	}, nil
}

func NewRoomWithName(
	name string,
	peers map[string]*Peer,
	tracks map[string]*Track,
) (*Room, error) {
	room, err := NewRoom(peers, tracks)
	if err != nil {
		return nil, err
	}

	room.Name = &name
	return room, nil
}
