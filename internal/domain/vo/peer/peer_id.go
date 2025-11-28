package peer

import "github.com/stazoloto/sfu-mediaserver/pkg/id"

type ID string

func NewPeerID() (ID, error) {
	peerID, err := id.GenerateIDWithPrefix("peer", 10)
	if err != nil {
		return "", err
	}
	return ID(peerID), nil
}
