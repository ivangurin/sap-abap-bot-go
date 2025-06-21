package bot

import (
	"context"
	"fmt"
)

func (s *Service) Run(ctx context.Context) error {
	me, err := s.bot.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("get bot info: %w", err)
	}

	s.username = me.Username

	s.bot.Start(ctx)

	return nil
}
