package entities

type MessageType string

const (
	TypeJoin      MessageType = "join"
	TypeLeave     MessageType = "leave"
	TypeOffer     MessageType = "offer"
	TypeAnswer    MessageType = "answer"
	TypeCandidate MessageType = "candidate"
	TypePeers     MessageType = "peers"
	TypeError     MessageType = "error"
)

type Message struct {
	Type     MessageType `json:"type"`
	Room     string      `json:"room,omitempty"`
	To       string      `json:"to,omitempty"`
	From     string      `json:"from,omitempty"`
	ClientID string      `json:"client_id,omitempty"`
	Payload  []byte      `json:"payload,omitempty"`
}
