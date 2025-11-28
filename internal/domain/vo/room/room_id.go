package room

import "github.com/stazoloto/sfu-mediaserver/pkg/id"

type ID string

func NewRoomID() (ID, error) {
	peerID, err := id.GenerateIDWithPrefix("room", 10)
	if err != nil {
		return "", err
	}
	return ID(peerID), nil
}
