package ws_app

import (
	"encoding/json"
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

func (s *wsGameService) Send(game *domain.Game, data []byte) error {
	return s.repoMap[game.PName].Send(data)
}

func (s *wsGameService) SendJSON(game *domain.Game, obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("SendJSON: %v", err)
	}
	return s.repoMap[game.PName].Send(data)
}

func (s *wsGameService) Receive(game *domain.Game) ([]byte, error) {
	return s.repoMap[game.PName].Receive()
}

func (s *wsGameService) Close(game *domain.Game) error {
	return s.repoMap[game.PName].Close()
}

func (s *wsGameService) SendStartConfirmation(game *domain.Game) error {
	return s.SendJSON(game, domain.MsgStart{Type: "start", GameId: game.Id, Color: game.Color})
}

func (s *wsGameService) SendToOpLocally(game *domain.Game, msg any) {
	opGame, ok := s.playerGameMap[game.OpName]
	if ok {
		opGame.LocalRecv <- msg
	}
}
