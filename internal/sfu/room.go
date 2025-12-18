package sfu

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type Room struct {
	mu     sync.RWMutex
	peers  map[string]*Peer
	tracks []*webrtc.TrackLocalStaticRTP
}

func NewRoom() *Room {
	return &Room{
		peers:  make(map[string]*Peer),
		tracks: make([]*webrtc.TrackLocalStaticRTP, 0),
	}
}

func (r *Room) AddPeer(p *Peer) {
	r.mu.Lock()
	r.peers[p.ClientID] = p
	r.mu.Unlock()
}

func (r *Room) RemovePeer(clientID string) {
	r.mu.Lock()
	delete(r.peers, clientID)
	r.mu.Unlock()
}

func (r *Room) GetPeer(clientID string) *Peer {
	r.mu.RLock()
	p := r.peers[clientID]
	r.mu.RUnlock()
	return p
}

func (r *Room) ForEachPeer(fn func(*Peer)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.peers {
		fn(p)
	}
}
