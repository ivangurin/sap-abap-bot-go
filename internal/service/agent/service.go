package agent

import (
	"context"

	github_client "bot/internal/client/github"
	pkg_config "bot/internal/config"
	"bot/internal/model"
	pkg_logger "bot/internal/pkg/logger"
)

type IService interface {
	ProcessPrompt(ctx context.Context, prompt string, threadMessages []*model.ThreadMessage) ([]*Answer, error)
}

type Service struct {
	config       *pkg_config.Config
	logger       *pkg_logger.Logger
	githubClient github_client.IClient
	tools        []*github_client.ChatCompletionRequestTool
}

func NewService(
	config *pkg_config.Config,
	logger *pkg_logger.Logger,
	githubClient github_client.IClient,
) IService {
	tools := []*github_client.ChatCompletionRequestTool{
		{
			Type: "function",
			Function: &github_client.ChatCompletionRequestFunction{
				Name:        "send_answer",
				Description: "Отправить ответ пользователю",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"answer": map[string]any{
							"type":        "string",
							"description": "Ответ пользователю",
						},
					},
				},
			},
		},
	}

	return &Service{
		config:       config,
		logger:       logger,
		githubClient: githubClient,
		tools:        tools,
	}
}
