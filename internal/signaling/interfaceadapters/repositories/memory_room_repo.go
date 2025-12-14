package repositories

import (
	"sync"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
)

type MemoryRoomRepo struct {
	mu    sync.RWMutex
	rooms map[string]*entities.Room
}

func NewMemoryRoomRepo() *MemoryRoomRepo {
	return &MemoryRoomRepo{
		mu:    sync.RWMutex{},
		rooms: make(map[string]*entities.Room),
	}
}

func (r *MemoryRoomRepo) GetOrCreate(id string) (*entities.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if room, ok := r.rooms[id]; ok {
		return room, nil
	}
	room := entities.NewRoom(id)
	r.rooms[id] = room
	return room, nil
}

func (r *MemoryRoomRepo) Get(id string) (*entities.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rooms[id], nil
}

func (r *MemoryRoomRepo) Save(room *entities.Room) error {
	return nil
}

func (r *MemoryRoomRepo) DeleteIfEmpty(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.rooms[id] != nil && len(r.rooms[id].Clients) == 0 {
		delete(r.rooms, id)
	}
	return nil
}
