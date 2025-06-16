package bot

func (s *Service) ErrorHandler(err error) {
	s.logger.Errorf("error handler: %s", err.Error())
}
