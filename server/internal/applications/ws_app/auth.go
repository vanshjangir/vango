package ws_app

import (
	"encoding/json"
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type wsAuthRequest struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (s *wsGameService) Auth(repo ports.WsGameRepository, us ports.UserService) (string, error) {
	data, err := repo.Receive()
	if err != nil {
		return "", err
	}
	var req wsAuthRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return "", fmt.Errorf("Error unmarhalling json")
	}

	switch req.Type {
	case "google":
		return us.AuthGoogle(req.Token)
	case "guest":
		return us.AuthGuest(req.Token)
	default:
		return "", fmt.Errorf("Unsupported token type")
	}
}
