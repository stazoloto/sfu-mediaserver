package sfu

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

type SFU struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	api   *webrtc.API

	// callbacks
	onICECandidate func(roomID, clientID string, c webrtc.ICECandidateInit)
	signalSender   func(to string, msg entities.Message)
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

func (s *SFU) SetSignalSender(fn func(to string, msg entities.Message)) {
	s.signalSender = fn
}

func (s *SFU) sendSignal(to string, msg entities.Message) {
	if s.signalSender != nil {
		s.signalSender(to, msg)
	}
}

func (s *SFU) SetOnICECandidate(fn func(roomID string, clientID string, c webrtc.ICECandidateInit)) {
	s.onICECandidate = fn
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

	_, _ = pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeVideo,
		webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		},
	)

	_, _ = pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		},
	)

	peer := &Peer{
		ClientID: clientID,
		PC:       pc,
	}

	room.AddPeer(peer)

	// колбэк
	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil || s.onICECandidate == nil {
			return
		}

		s.onICECandidate(
			roomID,
			clientID,
			c.ToJSON(),
		)
	})

	// remote - трек от клиента, local - трек, который сервер будет отдавать другим
	pc.OnTrack(func(remote *webrtc.TrackRemote, reciever *webrtc.RTPReceiver) {
		log.Printf(
			"SFU OnTrack: from=%s kind=%s id=%s stream=%s",
			clientID,
			remote.Kind(),
			remote.ID(),
			remote.StreamID(),
		)

		// создаем локальный трек
		local, err := webrtc.NewTrackLocalStaticRTP(
			remote.Codec().RTPCodecCapability,
			remote.ID(),
			remote.StreamID(),
		)
		if err != nil {
			return
		}

		// подключаем всем остальным пирам
		room.ForEachPeer(func(p *Peer) {
			if p.ClientID == clientID {
				return
			}
			sender, err := p.PC.AddTrack(local)
			if err != nil {
				return
			}

			go func() {
				rtcpBuf := make([]byte, 1500)
				for {
					if _, _, err := sender.Read(rtcpBuf); err != nil {
						return
					}
				}
			}()

			offer, err := p.PC.CreateOffer(nil)
			if err != nil {
				return
			}
			if err := p.PC.SetLocalDescription(offer); err != nil {
				return
			}

			b, _ := json.Marshal(offer)
			s.signalSender(p.ClientID, entities.Message{
				Type:    entities.TypeOffer,
				From:    "sfu",
				To:      p.ClientID,
				Room:    roomID,
				Payload: b,
			})
		})

		go func() {
			buf := make([]byte, 1500)
			for {
				n, _, err := remote.Read(buf)
				if err != nil {
					return
				}
				if _, err = local.Write(buf[:n]); err != nil {
					return
				}
			}
		}()
	})

	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state == webrtc.ICEConnectionState(webrtc.PeerConnectionStateDisconnected) || state == webrtc.ICEConnectionState(webrtc.PeerConnectionStateClosed) {
			room.RemovePeer(clientID)
			_ = pc.Close()
		}
	})

	return peer, nil
}

func (s *SFU) GetPeer(roomID, clientID string) *Peer {
	s.mu.RLock()
	room := s.rooms[roomID]
	s.mu.RUnlock()

	if room == nil {
		return nil
	}

	return room.GetPeer(clientID)
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
