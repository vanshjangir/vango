package user_app

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vanshjangir/vango/server/internal/applications/util_app"
)

func (s *userService) AuthGuest(token string) (string, error) {
	jwtToken, err := util_app.VerifyGuestToken(token)
	if err != nil {
		return "", fmt.Errorf("AuthGuest: VerifyGuestToken: %v", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("AuthGuest: Invalid jwt Claims")
	} else {
		return claims["username"].(string), nil
	}
}

func (s *userService) AuthGoogle(token string) (string, error) {
	tokenInfo, err := util_app.VerifyGoogleToken(token)
	if err != nil {
		return "", err
	}

	if user, err := s.ur.FindByEmail(tokenInfo.Email); err != nil {
		return "", err
	} else {
		return user.Username, nil
	}
}
