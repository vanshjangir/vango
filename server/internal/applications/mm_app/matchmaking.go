package mm_app

import (
	"fmt"
	"sync"
	"time"

	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

var (
	matched = make(map[string]bool)
	waiting = make(map[string]domain.WaitingPlayer)
	mu      sync.Mutex
	ttl     = 60 * time.Second
)

func addPlayer(p domain.WaitingPlayer) {
	mu.Lock()
	p.AddedAt = time.Now()
	waiting[p.Username] = p
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

func (s *matchMakingService) matchPlayer(p domain.WaitingPlayer) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	if ok := matched[p.Username]; ok {
		delete(matched, p.Username)
		return true, nil
	}

	now := time.Now()
	for id, w := range waiting {
		if now.Sub(w.AddedAt) > ttl {
			delete(waiting, id)
		}
	}

	for username, other := range waiting {
		if other.Username == p.Username {
			continue
		}
		if abs(p.Rating-other.Rating) <= 100 {
			twoWaitingPlayers := &[2]domain.WaitingPlayer{p, other}
			matched[p.Username] = true
			matched[username] = true

			if err := s.CreateNewGame(
				twoWaitingPlayers[0].Username,
				twoWaitingPlayers[1].Username,
			); err != nil {
				return false, fmt.Errorf("Match: %v", err)
			}

			delete(waiting, p.Username)
			delete(waiting, username)
			return true, nil
		}
	}
	return false, nil
}

func (s *matchMakingService) CreateNewGame(blackName, whiteName string) error {
	id, err := s.gr.CreateNewGame(blackName, whiteName)
	if err != nil {
		return fmt.Errorf("CreateNewGame: %v", err)
	}
	err = s.pr.SetGameWithUsername(id, blackName, whiteName)
	if err != nil {
		return fmt.Errorf("CreateNewGame: %v", err)
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
			removePlayer(wp.Username)
			return "", fmt.Errorf("Match: timeout")
		case <-tick:
			matched, err := s.matchPlayer(wp)
			if matched {
				wsurl := Pick()
				return wsurl, nil
			}
			if err != nil {
				return "", fmt.Errorf("Match: %v", err)
			}
		}
	}
}
