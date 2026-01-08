package game_app

import (
	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *gameService) RecentGames(username string, howmany int) ([]domain.GameReview, error) {
	game, err := s.gr.RecentGames(username, howmany)
	if err != nil {
		return []domain.GameReview{}, err
	}
	return game, err
}
