package agent

import (
	"bot/internal/model"
	"context"

	github_client "bot/internal/clients/github"
	pkg_config "bot/internal/config"

	pkg_logger "bot/internal/pkg/logger"
)

type Service interface {
	ProcessPrompt(ctx context.Context, prompt string, threadMessages []*model.ThreadMessage) ([]*Answer, error)
}

type service struct {
	config       *pkg_config.Config
	logger       pkg_logger.Logger
	githubClient github_client.Client
	tools        []*github_client.ChatCompletionRequestTool
}

func NewService(
	config *pkg_config.Config,
	logger pkg_logger.Logger,
	githubClient github_client.Client,
) Service {
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

	return &service{
		config:       config,
		logger:       logger,
		githubClient: githubClient,
		tools:        tools,
	}
}
