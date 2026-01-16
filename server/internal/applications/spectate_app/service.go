package spectate_app

import "github.com/vanshjangir/vango/server/internal/ports"

type spectateService struct {
	pr      ports.PubSubRepository
	repoMap map[int]map[ports.WsGameRepository]any
	msgChan chan string
}

func (s *spectateService) IsPlayerInGame(username string, gameId int) bool {
	game, err := s.pr.GetGameFromUsername(username)
	if err != nil {
		return false
	}
	if game.GameId == gameId {
		return true
	}
	return false
}

func NewSpectateService(pr ports.PubSubRepository) ports.SpectateService {
	return &spectateService{
		pr:      pr,
		repoMap: make(map[int]map[ports.WsGameRepository]any),
		msgChan: make(chan string),
	}
}
