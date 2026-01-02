package ports

import "github.com/vanshjangir/vango/server/internal/domain"

type GameService interface {
	SaveGame(g *domain.Game) error
	Review(gameid int) (domain.GameReview, error)
}

type UserService interface {
	ChangeUsername(oldName, newName string) error
	LoginByGoogle(credentials string) (string, error)
	LoginAsGuest(credentials string) (string, string, error)
	AuthGoogle(token string) (string, error)
	AuthGuest(token string) (string, error)
	GetUser(username string) (domain.User, error)
}

type WsGameService interface {
	Send(game *domain.Game, data []byte) error
	SendJSON(game *domain.Game, obj any) error
	Receive(game *domain.Game) ([]byte, error)
	Close(game *domain.Game) error
	Auth(authHeader string, us UserService) (string, error)
	Play(game *domain.Game)
	SetupGame(username string, repo WsGameRepository) (*domain.Game, error)
	LoadExistingGame(username string, repo WsGameRepository) (*domain.Game, error)
	SendStartConfirmation(game *domain.Game) error
	GetGameFromId(gameId int) *domain.Game
	GetGameFromPlayerName(name string) *domain.Game
}

type MatchMakingService interface {
	Match(domain.WaitingPlayer) (string, error)
}

type SpectateService interface {
	Start()
	AddSpectator(game *domain.Game, repo WsGameRepository)
	SendSyncState(blackGame, whiteGame *domain.Game, repo WsGameRepository)
}
