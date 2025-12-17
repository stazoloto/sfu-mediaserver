package usecase

import (
	"encoding/json"
	"errors"

	"github.com/pion/webrtc/v3"
	"github.com/stazoloto/sfu-mediaserver/internal/sfu"
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

	sfu *sfu.SFU
}

func NewInteractor(repo RoomRepository, out ClientGateway, sfu *sfu.SFU) *Interactor {
	i := &Interactor{
		rooms: repo,
		out:   out,
		sfu:   sfu,
	}

	sfu.SetOnICECandidate(i.handleSFUICE)
	sfu.SetSignalSender(func(to string, msg entities.Message) {
		_ = out.Send(to, msg)
	})
	return i
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

func (i *Interactor) HandleSFU(msg entities.Message) error {
	switch msg.Type {
	case entities.TypeOffer:
		return i.handleSFUOffer(msg)
	case entities.TypeCandidate:
		return i.handleSFUCandidate(msg)
	case entities.TypeAnswer:
		return i.handleSFUAnswer(msg)
	default:
		return nil
	}
}

func (i *Interactor) handleSFUOffer(msg entities.Message) error {
	if msg.Room == "" {
		return ErrMissingRoom
	}

	// парсинг SDP offer
	var offer webrtc.SessionDescription
	if err := json.Unmarshal(msg.Payload, &offer); err != nil {
		return err
	}

	// получить или созать peer в SFU
	peer, err := i.sfu.Join(msg.Room, msg.From)
	if err != nil {
		return err
	}

	// установить remote description
	if err := peer.PC.SetRemoteDescription(offer); err != nil {
		return err
	}

	// создание answer
	answer, err := peer.PC.CreateAnswer(nil)
	if err != nil {
		return err
	}

	// установить local description
	if err := peer.PC.SetLocalDescription(answer); err != nil {
		return err
	}

	// отправка answer клиенту
	answerBytes, _ := json.Marshal(answer)

	return i.out.Send(msg.From, entities.Message{
		Type:    entities.TypeAnswer,
		From:    "sfu",
		To:      msg.From,
		Room:    msg.Room,
		Payload: answerBytes,
	})
}

func (i *Interactor) handleSFUAnswer(msg entities.Message) error {
	if msg.Room == "" || msg.From == "" {
		return nil
	}
	peer := i.sfu.GetPeer(msg.Room, msg.From)
	if peer == nil {
		return nil
	}

	var answer webrtc.SessionDescription
	if err := json.Unmarshal(msg.Payload, &answer); err != nil {
		return err
	}

	return peer.PC.SetRemoteDescription(answer)
}

func (i *Interactor) handleSFUCandidate(msg entities.Message) error {
	peer := i.sfu.GetPeer(msg.Room, msg.From)
	if peer == nil {
		return nil
	}

	var c webrtc.ICECandidateInit
	if err := json.Unmarshal(msg.Payload, &c); err != nil {
		return err
	}

	_ = peer.PC.AddICECandidate(c)
	return nil
}

func (i *Interactor) handleSFUICE(roomID, clientID string, c webrtc.ICECandidateInit) {
	bytes, _ := json.Marshal(c)

	_ = i.out.Send(clientID, entities.Message{
		Type:    entities.TypeCandidate,
		From:    "sfu",
		To:      clientID,
		Room:    roomID,
		Payload: bytes,
	})
}

func (i *Interactor) handleSignal(to string) {

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
	if msg.Type == entities.TypeOffer {
		if msg.Room == "" || msg.From == "" {
			return nil
		}
	}

	if msg.Type == entities.TypeCandidate {
		if msg.Room == "" || msg.From == "" {
			return nil
		}
	}

	if msg.To == "sfu" {
		return i.HandleSFU(msg)
	}

	return i.out.Send(msg.To, msg)
}
