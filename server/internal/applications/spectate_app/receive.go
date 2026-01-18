package spectate_app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vanshjangir/vango/server/internal/domain"
	"github.com/vanshjangir/vango/server/internal/ports"
)

func (s *spectateService) AddSpectator(game *domain.Game, repo ports.WsGameRepository) {
	_, ok := s.repoMap[game.Id]
	if !ok {
		s.repoMap[game.Id] = make(map[ports.WsGameRepository]any)
	}
	s.repoMap[game.Id][repo] = struct{}{}
}

func (s *spectateService) SendToSpectators(payload domain.SpectateServicePayload) {
	for repo := range s.repoMap[payload.GameId] {
		data, err := json.Marshal(payload.GameData)
		if err != nil {
			log.Println("SendToSpectators: Marshal:", err)
			return
		}
		err = repo.Send(data)
		if err != nil {
			log.Println("Error sending data to a spectator:", err)
			repo.Close()
			delete(s.repoMap[payload.GameId], repo)
		}
	}
}

func (s *spectateService) SendStartConfirmation(repo ports.WsGameRepository) error {
	var startMsg domain.MsgStart
	startMsg.Type = "start"
	startMsg.Color = domain.BlackColor

	data, err := json.Marshal(startMsg)
	if err != nil {
		return fmt.Errorf("SendStartConfirmation: Marshal: %v", err)
	}
	
	err = repo.Send(data)
	if err != nil {
		return fmt.Errorf("SendStartConfirmation: repo.Send: %v", err)
	}
	return nil
}

func (s *spectateService) SendSyncState(blackGame, whiteGame *domain.Game, repo ports.WsGameRepository) error {
	var syncState domain.MsgSyncState
	syncState.Type = "syncstate"
	syncState.Gameid = blackGame.Id
	syncState.BlackName = blackGame.PName
	syncState.WhiteName = whiteGame.PName
	syncState.BlackRemTime = blackGame.GetRemainingTime()
	syncState.WhiteRemTime = whiteGame.GetRemainingTime()
	syncState.State, _ = blackGame.State.Board.Encode()
	syncState.History = blackGame.State.History

	data, err := json.Marshal(syncState)
	if err != nil {
		return fmt.Errorf("SendSyncState: Marshal: %v", err)
	}
	
	err = repo.Send(data)
	if err != nil {
		return fmt.Errorf("SendSyncState: repo.Send: %v", err)
	}

	return nil
}

func (s *spectateService) ReceiveGamesData() {
	go s.pr.Receive(s.msgChan)
	for {
		msg := <-s.msgChan
		var payload domain.SpectateServicePayload
		err := json.Unmarshal([]byte(msg), &payload)
		if err != nil {
			log.Println("Error unmarshaling json for SpectateServicePayload:", err)
			continue
		}

		go s.SendToSpectators(payload)
	}
}

func (s *spectateService) Start() {
	err := s.pr.Subscribe("spectator_game")
	if err != nil {
		log.Panicf("Error starting spectator service: %v\n", err)
	}
	go s.ReceiveGamesData()
}
