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
	return s.SendJSON(MsgStart{Type: "start"})
}

func (s *wsGameService) handleMove(game *domain.Game, data []byte) error {
	var msgMove MsgMove
	var msgMoveStatus MsgMoveStatus
	
	msgMoveStatus.Type = "movestatus"
	msgMoveStatus.Move = msgMove.Move
	msgMoveStatus.PassedTime = game.PassedTime
	
	err := json.Unmarshal(data, &msgMove)
	if err != nil {
		return fmt.Errorf("handleMove: unmarshalling msgMove: %v", err)
	}

	if (game.Turn != game.Color) {
		msgMoveStatus.InvalidTurn = true
		msgMoveStatus.InvalidMove = false
		data, err := json.Marshal(msgMoveStatus)
		if err != nil {
			return fmt.Errorf("handleMove: marshalling msgMoveStatus: %v",err)
		}
		if err := s.Send(data); err != nil {
			return fmt.Errorf("handleMove: Send invalid turn: %v", err)
		}
	}

	boardState, err := game.MakeMove(msgMove.Move)
	if err != nil {
		msgMoveStatus.InvalidTurn = false
		msgMoveStatus.InvalidMove = true
		data, err := json.Marshal(msgMoveStatus)
		if err != nil {
			return fmt.Errorf("handleMove: marshalling msgMoveStatus: %v",err)
		}
		if err := s.Send(data); err != nil {
			return fmt.Errorf("handleMove: Send invalid move: %v", err)
		}
	}

	msgMove.State = boardState
	msgMoveStatus.State = boardState
	
	s.sendToOpLocally(game, msgMove)
	if err := s.SendJSON(msgMoveStatus); err != nil {
		return fmt.Errorf("handleMove: sending msgMoveStatus: %v", err)
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

	game.Winner = winner
	game.WonBy = by
	
	s.sendToOpLocally(game, msgGameOver)
	if err := s.SendJSON(msgGameOver); err != nil {
		return err
	}
	
	return nil
}

func (s *wsGameService) handleAbort(game *domain.Game) error {
	err := s.handleGameOver(game, "abort", 1 - game.Color)
	if err != nil {
		return fmt.Errorf("handleAbort: handleGameover: %v", err)
	}
	return nil
}

func (s *wsGameService) handleLocalMsg(msg any) bool {
	switch msg.(type) {
	case MsgMove, MsgChat:
		if err := s.SendJSON(msg); err != nil {
			log.Println("handleLocalMsg: SendJSON for MsgMove, MsgChat:", err)
			return true
		}
	case MsgGameOver:
		if err := s.SendJSON(msg); err != nil {
			log.Println("handleLocalReceive: SendJSON for MsgGameOver", err)
		}
		return true
	}
	
	return false
}

func (s *wsGameService) ReceiveLocally(
	game *domain.Game,
	closeChan chan int,
){
	for {
		select {
		case <- closeChan:
			return
		case msg := <- game.LocalRecv:
			if s.handleLocalMsg(msg) == true {
				closeChan <- LOCAL_OUT;
				return
			}
		}
	}
}

func (s *wsGameService) handleChat(game *domain.Game, data []byte) error {
	var chatMsg MsgChat
	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return fmt.Errorf("handleChat: %v", err)
	}
	s.sendToOpLocally(game, chatMsg)
	return nil
}

func (s *wsGameService) handleClientData(game *domain.Game, data []byte) (bool, error) {
	shouldCancel := false
	var msgType MsgType
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

func (s *wsGameService) ReceiveFromClient(
	game *domain.Game,
	closeChan chan int,
	isOver *bool,
){
	for {
		data, err := s.Receive()
		if err != nil {
			if !(*isOver) {
				log.Println("Error receiving data", err)
				closeChan <- CLIENT_OUT
			}
			return
		}

		cancel, err := s.handleClientData(game, data)
		if err != nil {
			log.Println(err)
		}
		if cancel {
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

func (s *wsGameService) SaveGame(game *domain.Game) {
	if err := s.gr.SaveGame(game); err != nil {
		log.Println("Error in SaveGame:", err)
	}
}

func (s *wsGameService) Play(game *domain.Game) {
	closeChan := make(chan int)
	isOver := new(bool)
	*isOver = false
	go s.ReceiveFromClient(game, closeChan, isOver)
	go s.ReceiveLocally(game, closeChan)
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
	
	s.SaveGame(game);
}
