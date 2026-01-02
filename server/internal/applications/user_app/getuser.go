package user_app

import "github.com/vanshjangir/vango/server/internal/domain"

func (s *userService) GetUser(username string) (domain.User, error) {
	return s.ur.FindByUsername(username)
}
