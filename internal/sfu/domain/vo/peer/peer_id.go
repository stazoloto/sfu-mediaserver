package peer

import (
	"fmt"

	"github.com/stazoloto/sfu-mediaserver/pkg/id"
)

const (
	peerIDMin = 10000000
	peerIDMax = 99999999
)

type ID int64

func NewPeerID() (ID, error) {
	peerID, err := id.Generate(peerIDMin, peerIDMax)
	if err != nil {
		return 0, fmt.Errorf("generate peer ID: %w", err)
	}
	return ID(peerID), nil
}
