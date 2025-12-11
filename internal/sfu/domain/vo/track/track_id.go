package track

import (
	"fmt"

	"github.com/stazoloto/sfu-mediaserver/pkg/id"
)

const (
	sessionIDMin = 1000000000
	sessionIDMax = 9999999999
)

type ID int64

func NewTrackID() (ID, error) {
	trackId, err := id.Generate(sessionIDMin, sessionIDMax)

	if err != nil {
		return 0, fmt.Errorf("generate track id: %w", err)
	}

	return ID(trackId), nil
}
