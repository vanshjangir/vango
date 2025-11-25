package ws_app

import (
	"encoding/json"
	"fmt"
	"log"

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

	case "abort":
		err = s.handleAbort(game)
		shouldCancel = true

	case "chat":
		err = s.handleChat(game, data)

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
	var msgMoveStatus domain.MsgMoveStatus

	msgMoveStatus.Type = "movestatus"
	msgMoveStatus.Move = msgMove.Move
	msgMoveStatus.PassedTime = game.GetRemainingTime()
	msgMoveStatus.Code = "VALID"

	err := json.Unmarshal(data, &msgMove)
	if err != nil {
		return fmt.Errorf("handleMove: unmarshalling msgMove: %v", err)
	}

	if game.State.Turn != game.Color {
		msgMoveStatus.Code = "INVALID_TURN"
		err := s.SendJSON(game, msgMoveStatus);
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

	msgMove.State = boardState
	msgMoveStatus.State = boardState

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

	delete(s.gameMap, game.Id)
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
