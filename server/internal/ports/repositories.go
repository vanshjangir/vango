package ports

import (
	"github.com/vanshjangir/rapidgo/server/internal/domain"
)

type GameRepository interface {
	SaveGame(g *domain.Game) error
	Review(gameid int) (domain.GameReview, error)
	CreateNewGame(blackName, whiteName string) (int, error)
}

type UserRepository interface {
	ChangeUsername(oldName, newName string) error
	FindByUsername(username string) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	InsertUser(user domain.User) error
}

type WsGameRepository interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
}

type PubSubRepository interface {
	GetGameFromUsername(username string) (*domain.GameDataForPlayer, error)
	SetGameWithUsername(gameId int, blackName, whiteName string) error
}
