package spectate_app

import "github.com/vanshjangir/vango/server/internal/ports"

type spectateService struct {
	pr      ports.PubSubRepository
	repoMap map[int]map[ports.WsGameRepository]any
	msgChan chan string
}

func NewSpectateService(pr ports.PubSubRepository) ports.SpectateService {
	return &spectateService{
		pr:      pr,
		repoMap: make(map[int]map[ports.WsGameRepository]any),
		msgChan: make(chan string),
	}
}
