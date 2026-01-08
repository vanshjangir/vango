package user_app

import (
	"fmt"

	"github.com/vanshjangir/vango/server/internal/domain"
)

func (s *userService) GetUserProfile(username string) (domain.UserProfile, error) {
	var userProfile domain.UserProfile
	user, err := s.ur.FindByUsername(username)
	if err != nil {
		return userProfile, fmt.Errorf("FindByUsername: %v", err)
	}

	userProfile.Email = user.Email
	userProfile.Username = user.Username
	userProfile.Rating = user.Rating
	userProfile.Highestrating = user.Highestrating

	return userProfile, nil
}
