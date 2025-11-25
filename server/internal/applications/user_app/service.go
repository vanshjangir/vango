package user_app

import (
	"github.com/vanshjangir/vango/server/internal/ports"
)

type userService struct {
	ur ports.UserRepository
}

func NewUserService(ur ports.UserRepository) ports.UserService {
	return &userService{ur: ur}
}
