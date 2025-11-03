package ws_app

import (
	"github.com/vanshjangir/rapidgo/server/internal/domain"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type wsGameService struct {
	wr ports.WsGameRepository
	pr ports.PubSubRepository
    gr ports.GameRepository
    ur ports.UserRepository
	gameMap map[string]domain.Game
}

func NewWsGameService(
	wr ports.WsGameRepository,
	pr ports.PubSubRepository,
    gr ports.GameRepository,
    ur ports.UserRepository,
) ports.WsGameService {
	return &wsGameService{wr: wr, pr: pr, gr: gr, ur: ur}
}

func (s* wsGameService) CopyWsGameService(
	wr ports.WsGameRepository,
) ports.WsGameService {
	return &wsGameService{wr: wr, pr: s.pr, gr: s.gr, ur: s.ur}
}
