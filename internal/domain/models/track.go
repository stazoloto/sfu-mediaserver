package models

import (
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/room"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/track"
)

// Track - это медиапоток передаваемый от пользователя (аудио, видео, экран)
type Track struct {
	ID       track.ID        `json:"id"`
	PeerID   peer.ID         `json:"peer_id"`
	RoomID   room.ID         `json:"room_id"`
	Kind     track.TrackKind `json:"kind"`
	IsActive bool            `json:"is_active"`
}

func NewTrack(
	peerID peer.ID,
	roomID room.ID,
	kind track.TrackKind,
	isActive bool,
) *Track {
	id, err := track.NewTrackID()
	if err != nil {
		panic(err)
	}
	return &Track{
		ID:       id,
		PeerID:   peerID,
		RoomID:   roomID,
		Kind:     kind,
		IsActive: isActive,
	}
}
