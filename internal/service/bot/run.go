package bot

import (
	"context"
)

func (s *Service) Run(ctx context.Context) {
	s.bot.Start(ctx)
}
