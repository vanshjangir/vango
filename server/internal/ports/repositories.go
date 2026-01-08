package ports

import (
	"github.com/vanshjangir/vango/server/internal/domain"
)

type GameRepository interface {
	SaveGame(g *domain.Game) error
	Review(gameid int) (domain.GameReview, error)
	CreateNewGame(blackName, whiteName string) (int, error)
	RecentGames(username string, howmany int) ([]domain.GameReview, error)
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
	Send(game *domain.Game, msg any) error
	Subscribe(name string) error
	Unsubscribe(name string) error
	Receive(msgChan chan string)
}
