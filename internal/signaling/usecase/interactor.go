package usecase

import (
	"errors"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

var (
	ErrMissingRoom     = errors.New("missing room")
	ErrMissingClientID = errors.New("missing client")
	ErrMissingFrom     = errors.New("missing from")
	ErrMissingTo       = errors.New("missing to")
	ErrSamePeer        = errors.New("from and to must be different")
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
	if msg.Room == "" {
		return ErrMissingRoom
	}
	if msg.ClientID == "" {
		return ErrMissingClientID
	}

	room, err := i.rooms.GetOrCreate(msg.Room)
	if err != nil {
		return err
	}

	room.Clients[msg.ClientID] = &entities.Client{ID: msg.ClientID}
	_ = i.rooms.Save(room)

	peersPayload := serializePeers(room)

	for clientID := range room.Clients {
		i.out.Send(clientID, entities.Message{
			Type:    entities.TypePeers,
			Room:    room.ID,
			Payload: peersPayload,
		})

	}
	return nil
}

func (i *Interactor) Leave(msg entities.Message) error {
	if msg.Room == "" {
		return ErrMissingRoom
	}

	if msg.ClientID == "" {
		return ErrMissingClientID
	}
	room, err := i.rooms.GetOrCreate(msg.Room)
	if err != nil {
		return err
	}

	delete(room.Clients, msg.ClientID)
	_ = i.rooms.Save(room)
	_ = i.rooms.DeleteIfEmpty(room.ID)

	return i.out.Send(room.ID, entities.Message{
		Type:    entities.TypeLeave,
		Room:    room.ID,
		Payload: serializePeers(room),
	})
}

func (i *Interactor) Disconnect(clientID string) {
	rooms := i.rooms.GetAll()

	// Если клиента в комнате нет, пропуск
	for _, room := range rooms {
		if _, ok := room.Clients[clientID]; !ok {
			continue
		}

		delete(room.Clients, clientID)
		_ = i.rooms.Save(room)

		// комната опустела - удалить комнату
		if len(room.Clients) == 0 {
			_ = i.rooms.DeleteIfEmpty(room.ID)
			continue
		}

		peers := serializePeers(room)

		for peerID := range room.Clients {
			_ = i.out.Send(peerID, entities.Message{
				Type:    entities.TypePeers,
				Room:    room.ID,
				Payload: peers,
			})
		}
	}

}

func (i *Interactor) relay(msg entities.Message) error {
	if msg.From == "" {
		return ErrMissingFrom
	}
	if msg.To == "" {
		return ErrMissingTo
	}

	if msg.To == msg.From {
		return ErrSamePeer
	}
	return i.out.Send(msg.To, msg)
}
