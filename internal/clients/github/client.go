package github

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	pkg_config "bot/internal/config"
	"bot/internal/model"
	pkg_logger "bot/internal/pkg/logger"
)

type Client interface {
	Ask(ctx context.Context, systemPrompt string, messages []*model.ChatMessage) ([]string, error)
}

type client struct {
	config *pkg_config.Config
	logger pkg_logger.Logger
	openai openai.Client
}

func NewClient(
	config *pkg_config.Config,
	logger pkg_logger.Logger,
) Client {
	opts := []option.RequestOption{
		option.WithAPIKey(config.GitHub.Token),
		option.WithHeader("X-GitHub-Api-Version", "2022-11-28"),
	}
	if config.GitHub.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.GitHub.BaseURL))
	}

	return &client{
		config: config,
		logger: logger,
		openai: openai.NewClient(opts...),
	}
}
