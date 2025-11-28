package track

import (
	"github.com/stazoloto/sfu-mediaserver/pkg/id"
)

type ID string

func NewTrackID() (ID, error) {
	trackId, err := id.GenerateIDWithPrefix("track", 10)
	if err != nil {
		return "", err
	}
	return ID(trackId), nil
}
