package applications

import "fmt"

func (s *userService) LoginByEmail(email, password string) error {
	user, err := s.ur.FindByEmail(email)
	if err != nil {
		return err
	}
	if user.Password != email {
		return fmt.Errorf("Wrong password")
	}
	return nil
}
