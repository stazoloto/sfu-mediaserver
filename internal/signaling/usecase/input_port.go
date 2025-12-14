package usecase

import (
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

type InputPort interface {
	Handle(msg entities.Message) error
}
