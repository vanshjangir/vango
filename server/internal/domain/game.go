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
	BlackName string
	WhiteName string
	BlackTime int
	WhiteTime int
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

func (g *Game) Init(id, blackName, whiteName string, size, maxTime int) {
	g.Id = id
	g.BlackName = blackName
	g.WhiteName = whiteName
	g.Turn = WhiteTag
	g.BlackTime = 0
	g.WhiteTime = 0
	g.Size = size
	g.MaxTime = maxTime

	g.Board = new(baduk.Board)
	g.Board.Init(g.Size)
}

func (g *Game) MakeMove(move string, color int) (string, error) {
	if g.Turn == color {
		return "", fmt.Errorf("Not your turn")
	}

	nextturn := -1
	col := int(move[0] - 'a')
	row, err := strconv.Atoi(move[1:])
	if err != nil {
		return "", fmt.Errorf("Invalid move format")
	}

	if color == BlackTag {
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
	return g.Board.Encode()
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
