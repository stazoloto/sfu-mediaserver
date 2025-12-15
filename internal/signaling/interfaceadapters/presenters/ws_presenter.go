package presenters

import (
	"encoding/json"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

type Sender interface {
	Send(clientID string, data []byte) error
}

type WSPresenter struct {
	sender Sender
}

func NewWSPresenter(sender Sender) *WSPresenter {
	return &WSPresenter{sender: sender}
}

// -------------------------
// ClientGateway interface
// -------------------------

func (p *WSPresenter) Send(clientID string, msg entities.Message) error {
	data, err := p.encode(msg)
	if err != nil {
		return err
	}
	return p.sender.Send(clientID, data)
}

func (p *WSPresenter) encode(msg entities.Message) ([]byte, error) {
	dto := map[string]any{
		"type": msg.Type,
	}

	if msg.Room != "" {
		dto["room"] = msg.Room
	}
	if msg.From != "" {
		dto["from"] = msg.From
	}
	if msg.To != "" {
		dto["to"] = msg.To
	}
	if msg.ClientID != "" {
		dto["client_id"] = msg.ClientID
	}
	if len(msg.Payload) > 0 {
		dto["payload"] = json.RawMessage(msg.Payload)
	}

	return json.Marshal(dto)
}
