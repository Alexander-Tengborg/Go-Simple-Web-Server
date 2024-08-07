package poker

import "sync"

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{store: map[string]int{}}
}

type InMemoryPlayerStore struct {
	mu    sync.Mutex
	store map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mu.Lock()
	i.store[name]++
	i.mu.Unlock()
}

func (s *InMemoryPlayerStore) GetLeague() []Player {
	var players []Player

	for name, wins := range s.store {
		player := Player{name, wins}
		players = append(players, player)
	}

	return players
}
