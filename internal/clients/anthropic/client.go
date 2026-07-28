package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	pkg_config "bot/internal/config"
	"bot/internal/model"
	pkg_logger "bot/internal/pkg/logger"
)

type Client interface {
	Ask(ctx context.Context, systemPrompt string, messages []*model.ChatMessage) ([]string, error)
}

type client struct {
	config    *pkg_config.Config
	logger    pkg_logger.Logger
	anthropic anthropic.Client
}

func NewClient(
	config *pkg_config.Config,
	logger pkg_logger.Logger,
) Client {
	opts := []option.RequestOption{
		option.WithAPIKey(config.Anthropic.Token),
	}
	if config.Anthropic.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.Anthropic.BaseURL))
	}

	return &client{
		config:    config,
		logger:    logger,
		anthropic: anthropic.NewClient(opts...),
	}
}
