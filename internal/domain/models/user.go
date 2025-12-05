package models

import (
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
	"google.golang.org/grpc/peer"
)

// для авторизации НА БУДУЩЕЕ
type User struct {
	PeerID peer.ID
}
