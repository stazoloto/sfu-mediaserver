package usecase

import (
	"errors"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

type Interactor struct {
	rooms RoomRepository
	out   ClientGateway
}

func NewInteractor(repo RoomRepository, out ClientGateway) *Interactor {
	return &Interactor{
		rooms: repo,
		out:   out,
	}
}

func (i *Interactor) Handle(msg entities.Message) error {
	switch msg.Type {
	case entities.TypeJoin:
		return i.join(msg)
	case entities.TypeLeave:
		return i.Leave(msg)
	case entities.TypeAnswer, entities.TypeOffer, entities.TypeCandidate:
		return i.relay(msg)
	default:
		return errors.New("unknown message type")
	}
}

func (i *Interactor) join(msg entities.Message) error {
	room, err := i.rooms.GetOrCreate(msg.Room)
	if err != nil {
		return err
	}

	room.Clients[msg.ClientID] = &entities.Client{ID: msg.ClientID}
	_ = i.rooms.Save(room)

	return i.out.Broadcast(room.ID, entities.Message{
		Type:    entities.TypeJoin,
		Room:    room.ID,
		Payload: serializePeers(room)})
}

func (i *Interactor) Leave(msg entities.Message) error {
	room, err := i.rooms.GetOrCreate(msg.Room)
	if err != nil {
		return err
	}

	delete(room.Clients, msg.ClientID)
	_ = i.rooms.Save(room)
	_ = i.rooms.DeleteIfEmpty(room.ID)

	return i.out.Broadcast(room.ID, entities.Message{
		Type:    entities.TypeLeave,
		Room:    room.ID,
		Payload: serializePeers(room),
	})
}

func (i *Interactor) relay(msg entities.Message) error {
	if msg.To == "" {
		return errors.New("missing recipient")
	}
	return i.out.Broadcast(msg.To, msg)
}
