package usecase

import (
	"encoding/json"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

func serializePeers(room *entities.Room) []byte {
	ids := make([]string, 0, len(room.Clients))
	for id := range room.Clients {
		ids = append(ids, id)
	}

	b, _ := json.Marshal(ids)
	return b
}
