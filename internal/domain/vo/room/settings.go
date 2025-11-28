package room

import (
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
)

type Settings struct {
	MaxPeers             int       `json:"max_peers"`
	RecordingEnabled     bool      `json:"recording_enabled"`
	ScreenSharingEnabled bool      `json:"screen_sharing_enabled"`
	ChatEnabled          bool      `json:"chat_enabled"`
	DefaultRole          peer.Role `json:"default_role"`
	DefaultRoomType      RoomType  `json:"room_type"`
}
