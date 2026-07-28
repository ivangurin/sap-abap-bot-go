package agent

import (
	"bot/internal/model"
	"context"

	pkg_config "bot/internal/config"

	pkg_logger "bot/internal/pkg/logger"
)

type Service interface {
	ProcessPrompt(ctx context.Context, prompt string, threadMessages []*model.ThreadMessage) ([]*Answer, error)
}

// Client - выбранный ии провайдер, которому агент отправляет диалог.
type Client interface {
	Ask(ctx context.Context, systemPrompt string, messages []*model.ChatMessage) ([]string, error)
}

type service struct {
	config   *pkg_config.Config
	logger   pkg_logger.Logger
	aiClient Client
}

func NewService(
	config *pkg_config.Config,
	logger pkg_logger.Logger,
	aiClient Client,
) Service {
	return &service{
		config:   config,
		logger:   logger,
		aiClient: aiClient,
	}
}
