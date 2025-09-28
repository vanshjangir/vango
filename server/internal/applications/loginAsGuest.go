package applications

import (
	"strconv"

	"github.com/google/uuid"
)

func (s *userService) LoginAsGuest(credentials string) (string, string, error) {
	uniqueId := int(uuid.New().ID())
	username := "G" + strconv.Itoa(uniqueId)

	token, err := createGuestToken(username)
	if err != nil {
		return "", "", err
	}

	return token, username, nil
}
