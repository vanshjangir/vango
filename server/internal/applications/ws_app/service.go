package ws_app

import (
	"sync"

	"github.com/vanshjangir/vango/server/internal/domain"
	"github.com/vanshjangir/vango/server/internal/ports"
)

type wsGameService struct {
	pr      ports.PubSubRepository
	gr      ports.GameRepository
	ur      ports.UserRepository
	mu		*sync.Mutex
	gameMap map[int]*domain.Game
	playerGameMap map[string]*domain.Game
	repoMap map[string]ports.WsGameRepository
}

func safedelete[K comparable, V any](mu *sync.Mutex, m map[K]V, key K) {
    mu.Lock()
    delete(m, key)
    mu.Unlock()
}

func NewWsGameService(
	pr ports.PubSubRepository,
	gr ports.GameRepository,
	ur ports.UserRepository,
) ports.WsGameService {
	return &wsGameService{
		pr: pr, gr: gr, ur: ur,
		gameMap: make(map[int]*domain.Game),
		playerGameMap: make(map[string]*domain.Game),
		repoMap: make(map[string]ports.WsGameRepository),
		mu: new(sync.Mutex),
	}
}
