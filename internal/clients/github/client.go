package github

import (
	"context"

	pkg_config "bot/internal/config"
	pkg_logger "bot/internal/pkg/logger"
)

type Client interface {
	ChatCompletions(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error)
}

type client struct {
	config *pkg_config.Config
	logger pkg_logger.Logger
}

func NewClient(
	config *pkg_config.Config,
	logger pkg_logger.Logger,
) Client {
	return &client{
		config: config,
		logger: logger,
	}
}
