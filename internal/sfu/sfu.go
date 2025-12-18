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
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		return nil, err
	}

	// Сервер ожидает получать аудио/видео от клиента
	_, _ = pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeVideo,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	)
	_, _ = pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	)

	peer := &Peer{
		ClientID: clientID,
		PC:       pc,
		// ready=false по умолчанию, станет true после первого answer от клиента
	}

	// ---- callbacks ----

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil || s.onICECandidate == nil {
			return
		}
		s.onICECandidate(roomID, clientID, c.ToJSON())
	})

	pc.OnTrack(func(remote *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		s.handleOnTrack(roomID, room, peer, remote)
	})

	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state == webrtc.ICEConnectionStateDisconnected ||
			state == webrtc.ICEConnectionStateFailed ||
			state == webrtc.ICEConnectionStateClosed {
			room.RemovePeer(clientID)
			_ = pc.Close()
		}
	})

	// ---- register peer ONCE ----
	room.AddPeer(peer)

	// ---- subscribe new peer to already existing room tracks ----
	s.subscribePeerToExistingTracks(roomID, room, peer)

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

	if room, ok := s.rooms[roomID]; ok {
		return room
	}
	room := NewRoom()
	s.rooms[roomID] = room
	return room
}

func (s *SFU) handleOnTrack(roomID string, room *Room, from *Peer, remote *webrtc.TrackRemote) {
	log.Printf("SFU OnTrack: from=%s kind=%s", from.ClientID, remote.Kind())

	local, err := webrtc.NewTrackLocalStaticRTP(
		remote.Codec().RTPCodecCapability,
		remote.ID(),
		remote.StreamID(),
	)
	if err != nil {
		return
	}

	// сохранить трек в комнате (чтобы новые peers могли подписаться позже)
	room.mu.Lock()
	room.tracks = append(room.tracks, local)
	room.mu.Unlock()

	// раздать local всем остальным peers (НО RTP forward будет один раз, ниже)
	room.ForEachPeer(func(p *Peer) {
		if p.ClientID == from.ClientID {
			return
		}
		s.addTrackToPeer(roomID, p, local)
	})

	// один goroutine: remote -> local
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
}

func (s *SFU) subscribePeerToExistingTracks(roomID string, room *Room, peer *Peer) {
	room.mu.RLock()
	tracks := append([]*webrtc.TrackLocalStaticRTP(nil), room.tracks...)
	room.mu.RUnlock()

	for _, track := range tracks {
		s.addTrackToPeer(roomID, peer, track)
	}
}

func (s *SFU) addTrackToPeer(roomID string, peer *Peer, track *webrtc.TrackLocalStaticRTP) {
	sender, err := peer.PC.AddTrack(track)
	if err != nil {
		return
	}

	// читать RTCP обязательно
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, err := sender.Read(buf); err != nil {
				return
			}
		}
	}()

	s.requestRenegotiation(roomID, peer)
}

// requestRenegotiation шлёт offer только если можно.
// Если нельзя (не stable / уже negotiating) — ставит pending.
func (s *SFU) requestRenegotiation(roomID string, peer *Peer) {
	// Можно ли начать negotiation прямо сейчас?
	if !peer.BeginNegotiation() {
		peer.MarkPending()
		return
	}

	// Если не stable — нельзя, отложим
	if peer.PC.SignalingState() != webrtc.SignalingStateStable {
		peer.EndNegotiation()
		peer.MarkPending()
		return
	}

	offer, err := peer.PC.CreateOffer(nil)
	if err != nil {
		peer.EndNegotiation()
		return
	}

	if err := peer.PC.SetLocalDescription(offer); err != nil {
		peer.EndNegotiation()
		return
	}

	b, _ := json.Marshal(offer)
	s.sendSignal(peer.ClientID, entities.Message{
		Type:    entities.TypeOffer,
		From:    "sfu",
		To:      peer.ClientID,
		Room:    roomID,
		Payload: b,
	})
}

// OnAnswer должен вызываться из usecase, когда SFU получил answer от клиента.
// Он снимает negotiating, делает peer ready, и если был pending — запускает renegotiation ещё раз.
func (s *SFU) OnAnswer(roomID, clientID string) {
	peer := s.GetPeer(roomID, clientID)
	if peer == nil {
		return
	}

	needRetry := peer.OnAnswerApplied()
	if needRetry {
		// второй прогон negotiation (уже после того как peer стал ready и stable)
		s.requestRenegotiation(roomID, peer)
	}
}

func (s *SFU) RequestRenegotiation(roomID string, peer *Peer) {
	s.requestRenegotiation(roomID, peer)
}
