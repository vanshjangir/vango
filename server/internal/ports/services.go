package ports

import "github.com/vanshjangir/rapidgo/server/internal/domain"

type GameService interface {
	SaveGame(g *domain.Game) error
	Review(gameid int) (domain.GameReview, error)
}

type UserService interface {
	ChangeUsername(oldName, newName string) error
	LoginByGoogle(credentials string) (string, error)
	LoginAsGuest(credentials string) (string, string, error)
	Signup(username, email, password string) error
	AuthGoogle(token string) (string, error)
	AuthGuest(token string) (string, error)
}

type WsGameService interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
	Auth(us UserService) (string, error)
	Play(game *domain.Game)
	SetupGame(username string) (*domain.Game, error)
	SendStartConfirmation() error
    CopyWsGameService(wr WsGameRepository) WsGameService
}

type MatchMakingService interface {
	Match(domain.WaitingPlayer) (string, error)
}
