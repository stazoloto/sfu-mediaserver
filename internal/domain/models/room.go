package models

import "time"

type Room struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Peers     map[Peer.ID]*Peer
	Tracks    map[Track.ID]*Track
	MaxPeers  int       `json:"max_peers"`
	CreatedAt time.Time `json:"created_at"`
}
