package ws_app

import (
	"github.com/vanshjangir/rapidgo/server/internal/domain"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type wsGameService struct {
	wr ports.WsGameRepository
	pr ports.PubSubRepository
	gameMap map[string]domain.Game
}

func NewWsGameService(
	wr ports.WsGameRepository,
	pr ports.PubSubRepository,
) ports.WsGameService {
	return &wsGameService{wr: wr, pr: pr}
}

func (s *wsGameService) AddWsGameRepo(wr ports.WsGameRepository) {
	s.wr = wr
}
