package models

import (
	"time"

	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/room"
)

type Room struct {
	ID        room.ID           `json:"id"`
	Name      string            `json:"name"`
	Peers     map[string]*Peer  `json:"peers"`
	Tracks    map[string]*Track `json:"tracks"`
	Settings  room.Settings     `json:"settings"`
	CreatedAt time.Time         `json:"created_at"`
	IsActive  bool              `json:"is_active"`
}

func NewRoom(
	name string,
	peers map[string]*Peer,
	tracks map[string]*Track,
	settings room.Settings,
) *Room {
	id, err := room.NewRoomID()
	if err != nil {
		panic(err)
	}

	return &Room{
		ID:        id,
		Name:      name,
		Peers:     peers,
		Tracks:    tracks,
		Settings:  settings,
		CreatedAt: time.Now(),
		IsActive:  true,
	}
}
