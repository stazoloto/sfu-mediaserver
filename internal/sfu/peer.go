package sfu

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type Peer struct {
	ClientID string
	PC       *webrtc.PeerConnection

	mu          sync.Mutex
	ready       bool
	negotiating bool
	pending     bool
}

func (p *Peer) MarkReady() {
	p.mu.Lock()
	p.ready = true
	p.mu.Unlock()
}

// Ready значит: initial handshake завершён, peer может принимать renegotiation offers.
func (p *Peer) IsReady() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ready
}

// Пометить, что нужно renegotiation, но сейчас нельзя.
func (p *Peer) MarkPending() {
	p.mu.Lock()
	p.pending = true
	p.mu.Unlock()
}

// Забрать pending-флаг (и сбросить его).
func (p *Peer) TakePending() bool {
	p.mu.Lock()
	v := p.pending
	p.pending = false
	p.mu.Unlock()
	return v
}

// Пытаемся начать negotiation.
// Можно только если ready=true, negotiating=false.
func (p *Peer) BeginNegotiation() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.ready || p.negotiating {
		return false
	}
	p.negotiating = true
	return true
}

// Завершить negotiation (после успешного SetRemoteDescription(answer)).
func (p *Peer) EndNegotiation() {
	p.mu.Lock()
	p.negotiating = false
	p.mu.Unlock()
}

// Вызывается после применения answer.
// Делает peer ready (если ещё не был), снимает negotiating.
// Возвращает true, если есть pending renegotiation, которую надо добить.
func (p *Peer) OnAnswerApplied() (needRetry bool) {
	p.mu.Lock()
	p.ready = true
	p.negotiating = false
	needRetry = p.pending
	p.pending = false
	p.mu.Unlock()
	return
}
