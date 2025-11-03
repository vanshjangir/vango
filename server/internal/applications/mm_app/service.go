package mm_app

import (
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type matchMakingService struct {
	ur ports.UserRepository
	gr ports.GameRepository
	pr ports.PubSubRepository
}

func NewMatchMakingService(
	ur ports.UserRepository,
	gr ports.GameRepository,
	pr ports.PubSubRepository,
) *matchMakingService {
	return &matchMakingService{ur: ur, gr: gr, pr: pr}
}
