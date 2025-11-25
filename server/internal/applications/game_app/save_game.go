package game_app

import (
	"fmt"

	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *gameService) SaveGame(game *domain.Game) error {
	err := s.gr.SaveGame(game)
	if err != nil {
		return fmt.Errorf("SaveGame: %v", err)
	}
	return nil
}
