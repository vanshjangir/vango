package ws_app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

func (s *wsGameService) Send(data []byte) error {
	return s.wr.Send(data)
}

func (s *wsGameService) SendJSON(obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return s.wr.Send(data)
}

func (s *wsGameService) Receive() ([]byte, error) {
	return s.wr.Receive()
}

func (s *wsGameService) Close() error {
	return s.wr.Close()
}

func (s *wsGameService) SetupGame(username string) (*domain.Game, error) {
	gameData, err := s.pr.GetGameFromUsername(username)
	if err != nil {
		return nil, err
	}
	pname := gameData.BlackName
	opname := gameData.Whitename
	
	if (username != pname) {
		pname = gameData.Whitename
		opname = gameData.BlackName
	}

	game := new(domain.Game)
	game.Init(gameData.GameId, pname, opname, 19, 5*60*1000)

	return game, err
}

func (s *wsGameService) SendStartConfirmation() error {
	// send a proper message, MsgStart
	return s.Send([]byte("start"))
}

func (s *wsGameService) handleMove(game *domain.Game, data []byte) error {
	var msgMove MsgMove
	var msgMoveStatus MsgMoveStatus
	
	msgMoveStatus.Type = "movestatus"
	msgMoveStatus.Move = msgMove.Move
	msgMoveStatus.PassedTime = game.PassedTime
	
	err := json.Unmarshal(data, &msgMove)
	if err != nil {
		return err
	}

	if (game.Turn != game.Color) {
		msgMoveStatus.InvalidTurn = true
		msgMoveStatus.InvalidMove = false
		if data, err := json.Marshal(msgMoveStatus); err != nil {
			return err
		} else {
			s.Send(data)
		}
	}

	boardState, err := game.MakeMove(msgMove.Move)
	if err != nil {
		msgMoveStatus.InvalidTurn = false
		msgMoveStatus.InvalidMove = true
		if data, err := json.Marshal(msgMoveStatus); err != nil {
			return err
		} else {
			s.Send(data)
		}
	}

	msgMove.State = boardState
	msgMoveStatus.State = boardState
	
	s.sendToOpLocally(game, msgMove)
	if err := s.SendJSON(msgMoveStatus); err != nil {
		return err
	}
	
	return nil
}

func (s *wsGameService) sendToOpLocally(game *domain.Game, msg any) {
	opGame, ok := s.gameMap[game.OpName]
	if ok {
		opGame.LocalRecv <- msg
	}
}

func (s *wsGameService) handleGameOver(game *domain.Game, by string, winner int) error {
	var msgGameOver MsgGameOver
	msgGameOver.Type = "gameover"
	msgGameOver.By = by
	msgGameOver.Winner = winner
	
	s.sendToOpLocally(game, msgGameOver)
	if err := s.SendJSON(msgGameOver); err != nil {
		return err
	}
	
	return nil
}

func (s *wsGameService) handleAbort(game *domain.Game, closeChan chan int) error {
	err := s.handleGameOver(game, "abort", 1 - game.Color)
	if err != nil {
		return fmt.Errorf("handleGameover: %v", err)
	}
	closeChan <- CLIENT_OUT
	return nil
}

func (s *wsGameService) handleLocalReceive(
	game *domain.Game,
	closeChan chan int,
){
	for {
		select {
		case <- closeChan:
			return
		case <- game.LocalRecv:
			msg := <- game.LocalRecv
			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling json in handleLocalReceive", err)
				closeChan <- LOCAL_OUT
				return
			}

			err = s.Send(data)
			if err != nil {
				log.Println("Error sending data in handleLocalReceive", err)
				closeChan <- LOCAL_OUT
				return
			}
		}
	}
}

func (s *wsGameService) handleChat(game *domain.Game, data []byte) error {
	var chatMsg MsgChat
	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return err
	}
	s.sendToOpLocally(game, chatMsg)
	return nil
}

func (s *wsGameService) handleClientReceive(
	game *domain.Game,
	closeChan chan int,
	isOver *bool,
){
	for {
		data, err := s.Receive()
		if err != nil {
			if !*isOver {
				log.Println("Error receiving data", err)
				closeChan <- CLIENT_OUT
			}
			return
		}

		var msgType MsgType
		err = json.Unmarshal(data, &msgType)
		if err != nil {
			log.Println("Error unmarshlling json for MsgType", err)
			closeChan <- CLIENT_OUT
			return
		}

		switch msgType.Type {
		case "move":
			err = s.handleMove(game, data)
			err = fmt.Errorf("Error in handleMove: %v", err)
		case "abort":
			err = s.handleAbort(game, closeChan)
			err = fmt.Errorf("Error in handleAbort: %v", err)
		case "chat":
			err = s.handleChat(game, data)
			err = fmt.Errorf("Error in handleChat: %v", err)
		default:
			err = nil
		}

		if err != nil {
			log.Println(err)
			closeChan <- CLIENT_OUT
			return
		}
	}
}

func (s *wsGameService) checkTimer(game *domain.Game, closeChan chan int) {
	for {
		select {
		case <-closeChan:
			return
		default:
			if game.CheckTimeout() {
				closeChan <- TIMER_OUT
				return
			}
			time.Sleep(1*time.Second)
		}
	}
}

func (s *wsGameService) Play(game *domain.Game) {
	closeChan := make(chan int)
	isOver := new(bool)
	*isOver = false
	go s.handleClientReceive(game, closeChan, isOver)
	go s.handleLocalReceive(game, closeChan)
	go s.checkTimer(game, closeChan)
	
	out := <- closeChan
	switch (out) {
	case TIMER_OUT, LOCAL_OUT:
		*isOver = true
		err := s.Close()
		if err != nil {
			log.Println("Error closing websocket connection", err)
		}
	}

	close(closeChan)
}
