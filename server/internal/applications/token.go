package applications

import (
	"fmt"
	"os"
	"time"
	"context"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func verifyGuestToken(tokenString string) (*jwt.Token, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(
				"unexpected signing method: %v",
				token.Header["alg"],
			)
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func createGuestToken(username string) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(28 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyGoogleToken(token string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.NewService(
		context.Background(),
		option.WithAPIKey(os.Getenv("GOOGLE_CLIENT_ID")),
	)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(token).Do()
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}
