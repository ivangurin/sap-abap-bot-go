package bot

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) Close(ctx context.Context) error {
	_, err := s.bot.Close(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Too Many Requests") {
			return nil
		}
		return fmt.Errorf("close bot: %w", err)
	}

	return nil
}
