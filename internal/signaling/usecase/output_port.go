package usecase

import "github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"

type RoomRepository interface {
	GetOrCreate(roomID string) (*entities.Room, error)
	Get(roomID string) (*entities.Room, error)
	GetAll() []*entities.Room
	Save(room *entities.Room) error
	DeleteIfEmpty(roomID string) error
}

type ClientGateway interface {
	Send(clientID string, msg entities.Message) error
}
