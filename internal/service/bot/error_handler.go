package bot

func (s *service) ErrorHandler(err error) {
	s.logger.Errorf(s.ctx, "error handler: %s", err.Error())
}
