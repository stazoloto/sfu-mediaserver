package usecase

import "github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"

// Sent - хранит сообщения ("почтовый ящик")

type MockClientGateway struct {
	Sent map[string][]entities.Message
}

type BroadcastCall struct {
	RoomID string
	Msg    entities.Message
}

func NewMockClientGateway() *MockClientGateway {
	return &MockClientGateway{
		Sent: make(map[string][]entities.Message),
	}
}

func (m *MockClientGateway) Send(clientID string, msg entities.Message) error {
	m.Sent[clientID] = append(m.Sent[clientID], msg)
	return nil
}
