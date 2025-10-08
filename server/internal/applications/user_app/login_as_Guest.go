package user_app

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/vanshjangir/rapidgo/server/internal/applications/util_app"
)

func (s *userService) LoginAsGuest(credentials string) (string, string, error) {
	uniqueId := int(uuid.New().ID())
	username := "G" + strconv.Itoa(uniqueId)

	token, err := util_app.CreateGuestToken(username)
	if err != nil {
		return "", "", err
	}

	return token, username, nil
}
