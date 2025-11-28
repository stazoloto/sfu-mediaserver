package session

import "github.com/stazoloto/sfu-mediaserver/pkg/id"

type ID string

func NewSessionID() (ID, error) {
	trackId, err := id.GenerateIDWithPrefix("session", 10)
	if err != nil {
		return "", err
	}
	return ID(trackId), nil
}
