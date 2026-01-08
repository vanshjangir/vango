package ws_app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *wsGameService) GetGameFromId(gameId int) *domain.Game {
	return s.gameMap[gameId]
}

func (s *wsGameService) GetGameFromPlayerName(name string) *domain.Game {
	return s.playerGameMap[name]
}

func (s *wsGameService) handleClientData(game *domain.Game, data []byte) (bool, error) {
	shouldCancel := false
	var msgType domain.MsgType
	err := json.Unmarshal(data, &msgType)
	if err != nil {
		return true, fmt.Errorf("handleClientData: Unmarshal type: %v", err)
	}

	switch msgType.Type {
	case "move":
		err = s.handleMove(game, data)
		code := game.WinnerIfOver()
		if code != -1 {
			err = s.handleGameOver(game, "play", code)
			shouldCancel = true
		}

	case "abort":
		err = s.handleAbort(game)
		shouldCancel = true

	case "chat":
		err = s.handleChat(game, data)

	case "syncstate":
		err = s.handleSync(game)

	default:
		err = nil
	}

	if err != nil {
		err = fmt.Errorf("handleClientData: %v", err)
		shouldCancel = true
	}

	return shouldCancel, err
}

func (s *wsGameService) handleMove(game *domain.Game, data []byte) error {
	var msgMove domain.MsgMove
	err := json.Unmarshal(data, &msgMove)
	if err != nil {
		return fmt.Errorf("handleMove: unmarshalling msgMove: %v", err)
	}

	var msgMoveStatus domain.MsgMoveStatus
	msgMoveStatus.Type = "movestatus"
	msgMoveStatus.Code = "VALID"
	msgMoveStatus.Move = msgMove.Move
	if game.Color == domain.BlackColor {
		msgMoveStatus.BlackRemTime = game.GetRemainingTime()
		msgMoveStatus.WhiteRemTime = s.playerGameMap[game.OpName].GetRemainingTime()
	} else {
		msgMoveStatus.WhiteRemTime = game.GetRemainingTime()
		msgMoveStatus.BlackRemTime = s.playerGameMap[game.OpName].GetRemainingTime()
	}

	if game.State.Turn != game.Color {
		msgMoveStatus.Code = "INVALID_TURN"
		err := s.SendJSON(game, msgMoveStatus)
		if err != nil {
			return fmt.Errorf("handleMove: Send invalid turn: %v", err)
		}
		return nil
	}

	boardState, err := game.MakeMove(msgMove.Move)
	if err != nil {
		msgMoveStatus.Code = "INVALID_MOVE"
		if err := s.SendJSON(game, msgMoveStatus); err != nil {
			return fmt.Errorf("handleMove: Send invalid move: %v", err)
		}
		return nil
	}

	opGame := s.playerGameMap[game.OpName]
	if game.State.Turn != game.Color {
		game.RemTime -= int(time.Since(game.LastStoredTime).Milliseconds())
		opGame.LastStoredTime = time.Now()
	} else {
		opGame.RemTime -= int(time.Since(opGame.LastStoredTime).Milliseconds())
		game.LastStoredTime = time.Now()
	}

	msgMove.State = boardState
	msgMoveStatus.State = boardState
	msgMove.WhiteRemTime = msgMoveStatus.WhiteRemTime
	msgMove.BlackRemTime = msgMoveStatus.BlackRemTime

	s.SendToOpLocally(game, msgMove)

	err = s.SendJSON(game, msgMoveStatus)
	if err != nil {
		return fmt.Errorf("handleMove: sending msgMoveStatus: %v", err)
	}

	err = s.pr.Send(game, msgMove)
	if err != nil {
		return fmt.Errorf("handleMove: sending to pubsub: %v", err)
	}

	return nil
}

func (s *wsGameService) handleChat(game *domain.Game, data []byte) error {
	var chatMsg domain.MsgChat
	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return fmt.Errorf("handleChat: %v", err)
	}
	s.SendToOpLocally(game, chatMsg)
	return nil
}

func (s *wsGameService) handleGameOverWhenError(game *domain.Game, by string, winner int) error {
	var msgGameOver domain.MsgGameOver
	msgGameOver.Type = "gameover"
	msgGameOver.By = by
	msgGameOver.Winner = winner

	if err := s.SendJSON(game, msgGameOver); err != nil {
		return fmt.Errorf("handleGameOverWhenError: %v", err)
	}

	return nil
}

func (s *wsGameService) handleGameOver(game *domain.Game, by string, winner int) error {
	var msgGameOver domain.MsgGameOver
	msgGameOver.Type = "gameover"
	msgGameOver.By = by
	msgGameOver.Winner = winner

	game.Winner = winner
	game.WonBy = by

	s.SendToOpLocally(game, msgGameOver)
	if err := s.SendJSON(game, msgGameOver); err != nil {
		return err
	}

	return nil
}

func (s *wsGameService) handleAbort(game *domain.Game) error {
	err := s.handleGameOver(game, "abort", 1-game.Color)
	if err != nil {
		return fmt.Errorf("handleAbort: handleGameover: %v", err)
	}
	return nil
}

func (s *wsGameService) handleLocalMsg(game *domain.Game, msg any) bool {
	switch msg := msg.(type) {
	case domain.MsgMove, domain.MsgChat:
		if err := s.SendJSON(game, msg); err != nil {
			log.Println("handleLocalMsg: SendJSON for MsgMove, MsgChat:", err)
			return true
		}
	case domain.MsgGameOver:
		if err := s.SendJSON(game, msg); err != nil {
			log.Println("handleLocalReceive: SendJSON for MsgGameOver", err)
		}
		return true
	}

	return false
}

func (s *wsGameService) handleSync(game *domain.Game) error {
	var err error
	var msgSync domain.MsgSyncState
	msgSync.Type = "syncstate"
	msgSync.History = game.State.History
	msgSync.State, err = game.State.Board.Encode()
	if err != nil {
		return fmt.Errorf("handleSync: Encode: %v", err)
	}

	if game.Color == domain.WhiteColor {
		msgSync.WhiteName = game.PName
		msgSync.BlackName = game.OpName
		msgSync.WhiteRemTime = game.GetRemainingTime()
		msgSync.BlackRemTime = s.playerGameMap[game.OpName].GetRemainingTime()
	} else {
		msgSync.BlackName = game.PName
		msgSync.WhiteName = game.OpName
		msgSync.BlackRemTime = game.GetRemainingTime()
		msgSync.WhiteRemTime = s.playerGameMap[game.OpName].GetRemainingTime()
	}
	if err := s.SendJSON(game, msgSync); err != nil {
		return fmt.Errorf("handleSync: SendJSON: %v", err)
	}

	return nil
}
