package game_app

import (
	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *gameService) Review(gameid int) (domain.GameReview, error) {
	game, err := s.gr.Review(gameid)
	if err != nil {
		return domain.GameReview{}, err
	}
	return game, err
}
