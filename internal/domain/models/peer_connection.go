package models

import (
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/peer"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/room"
	"github.com/stazoloto/sfu-mediaserver/internal/domain/vo/track"
)

type PeerConnection struct {
	ID               int64                                    `json:"id"`
	PeerID           peer.ID                                  `json:"peer_id"`
	RoomID           room.ID                                  `json:"room_id"`
	Connection       *webrtc.PeerConnection                   `json:"-"` // Основное соединение
	PublishedTracks  map[track.ID]*webrtc.TrackRemote         ``         // Что пользователь публикует
	SubscribedTracks map[track.ID]*webrtc.TrackLocalStaticRTP // Что пользователь получает
	DataChannels     map[string]*webrtc.DataChannel           `json:"-"` // По ключу: "chat", "control", "sync"
	SignalingState   webrtc.SignalingState                    `json:"signaling_state"`
	ICEState         webrtc.ICEConnectionState                `json:"ice_state"`
	ConnectionState  webrtc.PeerConnectionState               `json:"connection_state"`
	GatheringState   webrtc.ICEGathererState                  `json:"gathering_state"`
	CreatedAt        time.Time                                `json:"created_at"`
	ConnectedAt      *time.Time                               `json:"connected_at,omitempty"`
	LastActivite     time.Time                                `json:"last_activity"`
}
