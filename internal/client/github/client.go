package github

import (
	"context"

	pkg_config "bot/internal/config"
	pkg_logger "bot/internal/pkg/logger"
)

type IClient interface {
	ChatCompletions(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error)
}

type ClientMock interface {
	IClient
}

type Client struct {
	config *pkg_config.Config
	logger *pkg_logger.Logger
}

func NewClient(
	config *pkg_config.Config,
	logger *pkg_logger.Logger,
) IClient {
	return &Client{
		config: config,
		logger: logger,
	}
}
