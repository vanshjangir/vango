package game_app

import (
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type gameService struct {
	gr ports.GameRepository
}

func NewGameService(gr ports.GameRepository) ports.GameService {
	return &gameService{gr: gr}
}
