package ws_app

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vanshjangir/rapidgo/server/internal/domain"
	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

func (s *wsGameService) SetupGame(username string, repo ports.WsGameRepository) (*domain.Game, error) {
	gameData, err := s.pr.GetGameFromUsername(username)
	if err != nil {
		return nil, err
	}
	pname := gameData.BlackName
	opname := gameData.Whitename

	if username != pname {
		pname = gameData.Whitename
		opname = gameData.BlackName
	}

	game := new(domain.Game)
	game.Init(
		gameData.GameId,
		pname, opname,
		19, 5*60*1000,
		//int(gameData.StartTime.Unix()),
		int(time.Now().Unix()), // for testing
	)
	if username == gameData.BlackName {
		game.Color = domain.BlackColor
	} else {
		game.Color = domain.WhiteColor
	}

	if _, ok := s.gameMap[game.OpName]; ok {
		game.State = s.gameMap[game.OpName].State
	} else {
		game.SetupState(19)
	}

	s.gameMap[username] = game
	s.repoMap[username] = repo
	return game, err
}

func (s *wsGameService) ReceiveLocally(game *domain.Game) {
	for {
		select {
		case <-game.CloseChan:
			return
		case msg := <-game.LocalRecv:
			if s.handleLocalMsg(game, msg) == true {
				game.CloseChan <- domain.GameCloseStatus{
					Code:           domain.LOCAL_OUT,
					ShouldSendToOp: false,
				}
				return
			}
		}
	}
}

func (s *wsGameService) ReceiveFromClient(game *domain.Game) {
	for {
		data, err := s.Receive(game)
		if err != nil {
			if !game.IsOver {
				game.IsOnline = false
				log.Println("Error receiving data", err)
				game.CloseChan <- domain.GameCloseStatus{
					Code:           domain.CLIENT_OUT,
					ShouldSendToOp: true,
				}
			}
			return
		}

		cancel, err := s.handleClientData(game, data)
		shouldSendToOp := false
		if err != nil {
			shouldSendToOp = true
			log.Println(err)
		}
		if cancel {
			game.CloseChan <- domain.GameCloseStatus{
				Code:           domain.CLIENT_OUT,
				ShouldSendToOp: shouldSendToOp,
			}
			return
		}
	}
}

func (s *wsGameService) checkTimer(game *domain.Game) {
	for {
		select {
		case <-game.CloseChan:
			return
		default:
			if game.CheckTimeout() {
				game.CloseChan <- domain.GameCloseStatus{
					Code:           domain.TIMER_OUT,
					ShouldSendToOp: false,
				}
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *wsGameService) LoadExistingGame(username string, repo ports.WsGameRepository) (*domain.Game, error) {
	game, ok := s.gameMap[username]
	if !ok {
		return nil, fmt.Errorf("Game not found")
	}
	s.repoMap[username] = repo
	return game, nil
}

func (s *wsGameService) SaveGame(game *domain.Game) {
	if err := s.gr.SaveGame(game); err != nil {
		log.Println("Error in SaveGame:", err)
	}
}

func (s *wsGameService) Play(game *domain.Game) {
	go s.ReceiveFromClient(game)
	go s.ReceiveLocally(game)
	go s.checkTimer(game)

	for {
		out := <-game.CloseChan
		if out.Code == domain.CLIENT_OUT && !game.IsOnline {
			select {
			case <-time.After(time.Second * domain.MAX_DISCN_TIME):
			case <-*game.ReconnectChan:
			}
			if game.IsOnline == true {
				go s.ReceiveFromClient(game)
				continue
			}
		}

		game.IsOver = true
		log.Println("Closing game...")

		if out.Code == domain.OP_INTERNAL_ERROR {
			if err := s.handleGameOverWhenError(game, "disconnection", game.Color); err != nil {
				log.Println("domain.INTERNAL_ERROR: ", err)
			}
		}

		if err := s.Close(game); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println("Error closing websocket connection:", err)
			}
		}

		if out.ShouldSendToOp {
			s.gameMap[game.OpName].IsOver = true
			s.gameMap[game.OpName].CloseChan <- domain.GameCloseStatus{
				Code:           domain.OP_INTERNAL_ERROR,
				ShouldSendToOp: false,
			}
		}

		close(game.CloseChan)
		delete(s.gameMap, game.PName)
		s.SaveGame(game)
		log.Println("Game Closed")
		break
	}
}
