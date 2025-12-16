package sfu

import (
	"errors"
	"sync"

	"github.com/pion/webrtc/v3"
)

type SFU struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	api   *webrtc.API
}

func NewSFU() *SFU {
	var m webrtc.MediaEngine
	_ = m.RegisterDefaultCodecs()

	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(&m),
	)

	return &SFU{
		rooms: map[string]*Room{},
		api:   api,
	}
}

func (s *SFU) Join(roomID, clientID string) (*Peer, error) {
	if roomID == "" || clientID == "" {
		return nil, errors.New("roomID and clientID required")
	}

	room := s.getOrCreateRoom(roomID)

	pc, err := s.api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		return nil, err
	}

	peer := &Peer{
		ClientID: clientID,
		PC:       pc,
	}

	room.AddPeer(peer)

	pc.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {})
}

func (s *SFU) getOrCreateRoom(roomID string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	// если комната есть - вернуть
	if room, ok := s.rooms[roomID]; ok {
		return room
	}

	// иначе создать
	room := NewRoom()
	s.rooms[roomID] = room
	return room
}
