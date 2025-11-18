package ws_app

import (
	"github.com/vanshjangir/rapidgo/server/internal/domain"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type wsGameService struct {
	pr      ports.PubSubRepository
	gr      ports.GameRepository
	ur      ports.UserRepository
	gameMap map[string]*domain.Game
	repoMap map[string]ports.WsGameRepository
}

func NewWsGameService(
	pr ports.PubSubRepository,
	gr ports.GameRepository,
	ur ports.UserRepository,
) ports.WsGameService {
	return &wsGameService{
		pr: pr, gr: gr, ur: ur,
		gameMap: make(map[string]*domain.Game),
		repoMap: make(map[string]ports.WsGameRepository),
	}
}
