package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vanshjangir/baduk"
)

const (
	WhiteColor = 0
	BlackColor = 1
	Handicap = 7.5
)

const (
	TIMER_OUT = 1
	CLIENT_OUT = 2
	LOCAL_OUT = 3
	INTERNAL_ERROR = 4
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

type Game struct {
	Id        int
	Size      int
	MaxTime   int
	Board     *baduk.Board
	Color	  int
	PName     string
	OpName 	  string
	LocalRecv chan any
	PassedTime int
    Winner     int
    WonBy      string
	Turn      int
	History   StringArray

	CloseChan chan GameCloseStatus
	IsOver    bool
}

type GameReview struct {
	Id        int
	BlackName string
	WhiteName string
	Winner    int
	Moves   StringArray
}

type GameCloseStatus struct {
	Code	int
	ShouldSendToOp bool
}

func (g *Game) Init(id int, playerName, opponentName string, size, maxTime int) {
	g.Id = id
	g.PName = playerName
	g.OpName = opponentName
	g.Turn = BlackColor
	g.PassedTime = 0
	g.Size = size
	g.MaxTime = maxTime

	g.LocalRecv = make(chan any)

	g.Board = new(baduk.Board)
	g.Board.Init(g.Size)
}

func (g *Game) MakeMove(move string) (string, error) {
	col := int(move[0] - 'a')
	row, err := strconv.Atoi(move[1:])
	if err != nil {
		return "", fmt.Errorf("Invalid move format")
	}

	if g.Color == BlackColor {
		err = g.Board.SetB(col, row)
	} else {
		err = g.Board.SetW(col, row)
	}
	if err != nil {
		return "", err
	}

	g.Turn = 1 - g.Turn
	g.History = append(g.History, move)
	encode, err := g.Board.Encode()
	return encode, err
}

func (g *Game) CheckTimeout() bool {
    if g.Turn != g.Color {
        return false;
    }

    if g.PassedTime >= g.MaxTime {
        return true
    }

    return false
}

func (g *Game) WinnerIfOver() int {
	total := len(g.History)
	if total < 2 {
		return -1
	}

	doublePass := (g.History[total-1] == "ps") && (g.History[total-2] == "ps")
	numPieces := 0
	for i := range g.Board.Grid {
		for j := range g.Board.Grid[i] {
			if g.Board.Grid[i][j].Empty {
				numPieces += 1
			}
		}
	}

	if doublePass || (numPieces == g.Size*g.Size) {
		bs, ws := g.Board.Score()
		if float32(bs) > float32(ws)+Handicap {
			return BlackColor
		} else {
			return WhiteColor
		}
	}

	return -1
}
