package mm_app

import (
	"fmt"
	"sync"
	"time"

	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

var (
	waiting = make(map[string]domain.WaitingPlayer)
	mu      sync.Mutex
	ttl     = 60 * time.Second
)

func addPlayer(p domain.WaitingPlayer) {
	mu.Lock()
	p.AddedAt = time.Now()
	waiting[p.Id] = p
	mu.Unlock()
}

func removePlayer(playerID string) {
	mu.Lock()
	delete(waiting, playerID)
	mu.Unlock()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func matchPlayer(p domain.WaitingPlayer) *[2]domain.WaitingPlayer {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for id, w := range waiting {
		if now.Sub(w.AddedAt) > ttl {
			delete(waiting, id)
		}
	}

	for id, other := range waiting {
		if other.Id == p.Id {
			continue
		}
		if abs(p.Rating-other.Rating) <= 100 {
			delete(waiting, p.Id)
			delete(waiting, id)
			return &[2]domain.WaitingPlayer{p, other}
		}
	}
	return nil
}

func (s *matchMakingService) Match(wp domain.WaitingPlayer) (string, error) {
	addPlayer(wp)

	timeout := time.After(30 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			removePlayer(wp.Id)
			return "", fmt.Errorf("timeout")
		case <-tick:
			if mp := matchPlayer(wp); mp != nil {
                s.gr.CreateNewGame(mp[0].Username, mp[1].Username)
                wsurl := Pick();
                return wsurl, nil
			}
		}
	}
}
