package room

import (
	"fmt"

	"github.com/stazoloto/sfu-mediaserver/pkg/id"
)

const (
	RoomIDMin = 10000000000000
	RoomIDMax = 99999999999999
)

type ID int64

func NewRoomID() (ID, error) {
	roomID, err := id.Generate(RoomIDMin, RoomIDMax)
	if err != nil {
		return 0, fmt.Errorf("generate room ID: %w", err)
	}

	return ID(roomID), nil
}
