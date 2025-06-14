package bot

func (s *Service) ErrorHandler(err error) {
	s.logger.Errorf("Error in bot handler: %s", err.Error())
}
