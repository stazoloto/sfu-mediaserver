package controllers

import (
	"encoding/json"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/usecase"
)

type WSMessageDTO struct {
	Type     string          `json:"type"`
	Room     string          `json:"room,omitempty"`
	From     string          `json:"from,omitempty"`
	To       string          `json:"to,omitempty"`
	ClientID string          `json:"client_id,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

type WSController struct {
	input usecase.InputPort
}

func NewWSController(input usecase.InputPort) *WSController {
	return &WSController{input: input}
}

func (c *WSController) Handle(raw []byte) error {
	var dto WSMessageDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return err
	}

	return c.input.Handle(entities.Message{
		Type:     entities.MessageType(dto.Type),
		Room:     dto.Room,
		From:     dto.From,
		To:       dto.To,
		ClientID: dto.ClientID,
		Payload:  dto.Payload})
}
