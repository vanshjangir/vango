package game_app

import (
	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

func (s *gameService) Review(gameid string) (domain.GameReview, error) {
	game, err := s.gr.Review(gameid)
	if err != nil {
		return domain.GameReview{}, err
	}
	return game, err
}
