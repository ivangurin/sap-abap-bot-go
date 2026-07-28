package github

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"

	"bot/internal/model"
)

const temperature = 1

// Ask отправляет диалог в GitHub Models (OpenAI-совместимый API) и возвращает ответ пользователю.
func (c *client) Ask(ctx context.Context, systemPrompt string, messages []*model.ChatMessage) ([]string, error) {
	requestMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)
	requestMessages = append(requestMessages, openai.SystemMessage(systemPrompt))

	for _, message := range messages {
		if message.Role == model.ChatRoleAssistant {
			requestMessages = append(requestMessages, openai.AssistantMessage(message.Text))
			continue
		}
		requestMessages = append(requestMessages, openai.UserMessage(message.Text))
	}

	if c.logger.IsWithDebug() {
		c.logger.Debugf(ctx, "github request: model=%s, messages=%d", c.config.GitHub.AIModel, len(requestMessages))
	}

	resp, err := c.openai.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       c.config.GitHub.AIModel,
		Messages:    requestMessages,
		Temperature: openai.Float(temperature),
	})
	if err != nil {
		return nil, fmt.Errorf("chat completions: %w", err)
	}

	answers := []string{}
	for _, choice := range resp.Choices {
		if choice.Message.Content == "" {
			continue
		}
		answers = append(answers, choice.Message.Content)
	}

	return answers, nil
}
