package user_app

import (
	"database/sql"
	"fmt"

	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *userService) Signup(username, email, password string) error {
	_, err := s.ur.FindByEmail(email)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return fmt.Errorf("Email already in use")
	}

	_, err = s.ur.FindByUsername(username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return fmt.Errorf("Username already in use")
	}

	err = s.ur.InsertUser(domain.User{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}
