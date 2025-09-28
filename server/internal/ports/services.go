package ports

import "github.com/vanshjangir/rapidgo/server/internal/domain"

type GameService interface {
	SaveGame(g *domain.Game) error
	Review(gameid string) (domain.GameReview, error)
}

type UserService interface {
	ChangeUsername(oldName, newName string) error
	LoginByEmail(email, password string) error
	LoginByGoogle(credentials string) (string, error)
	LoginAsGuest(credentials string) (string, string, error)
	Signup(username, email, password string) error
	AuthGoogle(token string) (string, error)
	AuthGuest(token string) (string, error)
	Match(domain.WaitingPlayer) (string, error)
}

type WsGameService interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
	Auth(us UserService) (string, error)
	Play(username string)
}
