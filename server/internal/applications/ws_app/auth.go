package ws_app

import (
	"fmt"
	"strings"

	"github.com/vanshjangir/vango/server/internal/ports"
)

func (s *wsGameService) Auth(authHeader string, us ports.UserService) (string, error) {
	tokenData := strings.Split(authHeader, " ")
	tokenType := tokenData[0]
	token := tokenData[1]

	switch tokenType {
	case "google":
		return us.AuthGoogle(token)
	case "guest":
		return us.AuthGuest(token)
	default:
		return "", fmt.Errorf("Unsupported token type")
	}
}
