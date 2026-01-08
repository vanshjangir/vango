package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/vanshjangir/baduk"
)

const (
	WhiteColor = 0
	BlackColor = 1
	Handicap   = 7.5

	TIMER_OUT         = "timer_out"
	CLIENT_OUT        = "client_out"
	LOCAL_OUT         = "local_out"
	OP_INTERNAL_ERROR = "op_internal_error"

	MAX_DISCN_TIME = 30
)

type StringArray []string

func (a *StringArray) Scan(src any) error {
	if src == nil {
		*a = []string{}
		return nil
	}
	return json.Unmarshal(src.([]byte), a)
}

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type GameState struct {
	History StringArray
	Size    int
	Board   *baduk.Board
	Turn    int
}

type Game struct {
	Id             int
	RemTime        int
	LastStoredTime time.Time
	Color          int
	PName          string
	OpName         string
	LocalRecv      chan any
	Winner         int
	WonBy          string

	CloseChan     chan GameCloseStatus
	ReconnectChan *chan bool
	IsOver        bool
	IsOnline      bool

	State *GameState
}

type GameReview struct {
	Gameid    int
	BlackName string
	WhiteName string
	Winner    int
	Moves     StringArray
	CreatedAt time.Time
}

type GameCloseStatus struct {
	Code           string
	ShouldSendToOp bool
}

func (g *Game) SetupState(size int) {
	g.State = new(GameState)
	g.State.Turn = BlackColor
	g.State.Board = new(baduk.Board)
	g.State.Size = size
	g.State.Board.Init(g.State.Size)
}

func (g *Game) Init(id int, playerName, opponentName string, size, maxTime int) {
	g.Id = id
	g.PName = playerName
	g.OpName = opponentName
	g.RemTime = maxTime
	g.LastStoredTime = time.Now()
	g.CloseChan = make(chan GameCloseStatus)
	g.IsOver = false
	g.IsOnline = true

	reconnectChan := make(chan bool)
	g.ReconnectChan = &reconnectChan

	g.LocalRecv = make(chan any)
}

func (g *Game) GetRemainingTime() int {
	if g.Color == g.State.Turn {
		return g.RemTime - int(time.Since(g.LastStoredTime).Milliseconds())
	} else {
		return g.RemTime
	}
}

func (g *Game) MakeMove(move string) (string, error) {
	if move == "ps" {
		g.State.Turn = 1 - g.State.Turn
		g.State.History = append(g.State.History, move)
		encode, err := g.State.Board.Encode()
		return encode, err
	}

	col := int(move[0] - 'a')
	row, err := strconv.Atoi(move[1:])
	if err != nil {
		return "", fmt.Errorf("Invalid move format")
	}

	if g.Color == BlackColor {
		err = g.State.Board.SetB(col, row)
	} else {
		err = g.State.Board.SetW(col, row)
	}
	if err != nil {
		return "", err
	}

	g.State.Turn = 1 - g.State.Turn
	g.State.History = append(g.State.History, move)
	encode, err := g.State.Board.Encode()
	return encode, err
}

func (g *Game) CheckTimeout() bool {
	if g.State.Turn != g.Color {
		return false
	}

	if g.GetRemainingTime() > 0 {
		return false
	}

	return true
}

func (g *Game) WinnerIfOver() int {
	total := len(g.State.History)
	if total < 2 {
		return -1
	}

	doublePass := (g.State.History[total-1] == "ps") && (g.State.History[total-2] == "ps")
	numPieces := 0
	for i := range g.State.Board.Grid {
		for j := range g.State.Board.Grid[i] {
			if g.State.Board.Grid[i][j].Empty {
				numPieces += 1
			}
		}
	}

	if doublePass || (numPieces == g.State.Size*g.State.Size) {
		bs, ws := g.State.Board.Score()
		if float32(bs) > float32(ws)+Handicap {
			return BlackColor
		} else {
			return WhiteColor
		}
	}

	return -1
}
