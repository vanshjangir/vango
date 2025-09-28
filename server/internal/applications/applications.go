package applications

import (
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type gameService struct {
	gr ports.GameRepository
}

type userService struct {
	ur ports.UserRepository
}

type wsGameService struct {
	wr ports.WsGameRepository
}

func NewGameService(gr ports.GameRepository) ports.GameService {
	return &gameService{gr: gr}
}

func NewUserService(ur ports.UserRepository) ports.UserService {
	return &userService{ur: ur}
}

func NewWsGameService(wr ports.WsGameRepository) ports.WsGameService {
    return &wsGameService{wr: wr}
}
