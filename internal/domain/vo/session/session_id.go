package session

import (
	"fmt"

	"github.com/stazoloto/sfu-mediaserver/pkg/id"
)

const (
	SessionIDMin = 1000000000
	SessionIDMax = 9999999999
)

type ID int64

func NewSessionID() (ID, error) {
	sessionId, err := id.Generate(SessionIDMin, SessionIDMax)
	if err != nil {
		return 0, fmt.Errorf("generate session id: %w", err)
	}
	return ID(sessionId), nil
}
