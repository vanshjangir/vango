package domain

import (
	"fmt"
	"strconv"

	"github.com/vanshjangir/baduk"
)

const (
	WhiteTag = 0
	BlackTag = 1
	Handicap = 7.5
)

type Game struct {
	Id        string
	Size      int
	MaxTime   int
	Board     *baduk.Board
	Color	  int
	PName     string
	OpName 	  string
	LocalRecv chan any
	PassedTime int
	Turn      int
	History   []string
}

type GameReview struct {
	Id        string
	BlackName string
	WhiteName string
	Winner    int
	History   []string
}

func (g *Game) Init(id, playerName, opponentName string, size, maxTime int) {
	g.Id = id
	g.PName = playerName
	g.OpName = opponentName
	g.Turn = BlackTag
	g.PassedTime = 0
	g.Size = size
	g.MaxTime = maxTime

	g.LocalRecv = make(chan any)

	g.Board = new(baduk.Board)
	g.Board.Init(g.Size)
}

func (g *Game) MakeMove(move string) (string, error) {
	nextturn := -1
	col := int(move[0] - 'a')
	row, err := strconv.Atoi(move[1:])
	if err != nil {
		return "", fmt.Errorf("Invalid move format")
	}

	if g.Color == BlackTag {
		err = g.Board.SetB(col, row)
		nextturn = WhiteTag
	} else {
		err = g.Board.SetW(col, row)
		nextturn = BlackTag
	}
	if err != nil {
		return "", err
	}

	g.Turn = nextturn
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
			return BlackTag
		} else {
			return WhiteTag
		}
	}

	return -1
}
