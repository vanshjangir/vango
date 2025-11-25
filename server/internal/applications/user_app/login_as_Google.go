package user_app

import (
	"database/sql"
	"strconv"

	"github.com/google/uuid"
	"github.com/vanshjangir/vango/server/internal/applications/util_app"
	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *userService) LoginByGoogle(credentials string) (string, error) {
	tokenInfo, err := util_app.VerifyGoogleToken(credentials)
	if err != nil {
		return "", err
	}

	var user domain.User
	user, err = s.ur.FindByEmail(tokenInfo.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			uniqueId := int(uuid.New().ID())
			username := "U" + strconv.Itoa(uniqueId)
			err = s.ur.InsertUser(domain.User{
				Username: username,
				Email:    tokenInfo.Email,
			})
			if err != nil {
				return "", err
			} else {
				return username, nil
			}
		} else {
			return "", err
		}
	}
	return user.Username, nil
}
