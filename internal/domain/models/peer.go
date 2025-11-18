package models

type Peer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Session
}
