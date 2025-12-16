package sfu

import "github.com/pion/webrtc/v3"

type Peer struct {
	ClientID string
	PC       *webrtc.PeerConnection
}
