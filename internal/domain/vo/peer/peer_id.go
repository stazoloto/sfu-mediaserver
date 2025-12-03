package peer

import (
	"fmt"

	"github.com/stazoloto/sfu-mediaserver/pkg/id"
)

const (
	PeerIDMin = 10000000
	PeerIDMax = 99999999
)

type ID int64

func NewPeerID() (ID, error) {
	peerID, err := id.Generate(PeerIDMin, PeerIDMax)
	if err != nil {
		return 0, fmt.Errorf("generate peer ID: %w", err)
	}
	return ID(peerID), nil
}
