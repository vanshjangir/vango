package applications

import (
	"encoding/json"
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/ports"
)

type WsAuthRequest struct {
	Type	string	`json:"type"`
	Token	string	`json:"token"`
}

func (s *wsGameService) Send(data []byte) error {
	return s.wr.Send(data)
}

func (s *wsGameService) Receive() ([]byte, error) {
	return s.wr.Receive()
}

func (s *wsGameService) Close() error {
	return s.wr.Close()
}

func (s *wsGameService) Auth(us ports.UserService) (string, error) {
	data, err := s.wr.Receive()
	if err != nil {
		return "", err
	}
	var req WsAuthRequest
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

func (s *wsGameService) Play(username string) {
}
