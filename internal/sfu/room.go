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
		peers: make(map[string]*Peer),
	}
}

func (r *Room) AddPeer(p *Peer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.peers[p.ClientID] = p
}

func (r *Room) RemovePeer(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.peers, id)
}

func (r *Room) ForEachPeer(fn func(*Peer)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.peers {
		fn(p)
	}
}

func (r *Room) GetPeer(clientID string) *Peer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.peers[clientID]
}
