package applications

func (s *userService) ChangeUsername(oldName string, newName string) error {
	err := s.ur.ChangeUsername(oldName, newName)
	if err != nil {
		return err
	}
	return nil
}
