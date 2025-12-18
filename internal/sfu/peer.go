package sfu

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type Peer struct {
	ClientID string
	PC       *webrtc.PeerConnection

	negotiating bool
	mu          sync.Mutex
}

func (p *Peer) BeginNegotiation() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.negotiating {
		return false
	}

	p.negotiating = true
	return true
}

func (p *Peer) EndNegotiation() {
	p.mu.Lock()
	p.negotiating = false
	p.mu.Unlock()
}
