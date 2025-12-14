package entities

type Client struct {
	ID string
}

type Room struct {
	ID      string
	Clients map[string]*Client
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Clients: make(map[string]*Client),
	}
}
