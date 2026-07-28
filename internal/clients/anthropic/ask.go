package anthropic

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"

	"bot/internal/model"
)

const maxTokens = 4096

// Ask отправляет диалог в Anthropic и возвращает ответ пользователю.
func (c *client) Ask(ctx context.Context, systemPrompt string, messages []*model.ChatMessage) ([]string, error) {
	requestMessages := make([]anthropic.MessageParam, 0, len(messages))
	for _, message := range messages {
		block := anthropic.NewTextBlock(message.Text)
		if message.Role == model.ChatRoleAssistant {
			requestMessages = append(requestMessages, anthropic.NewAssistantMessage(block))
			continue
		}
		requestMessages = append(requestMessages, anthropic.NewUserMessage(block))
	}

	if c.logger.IsWithDebug() {
		c.logger.Debugf(ctx, "anthropic request: model=%s, messages=%d", c.config.Anthropic.AIModel, len(requestMessages))
	}

	resp, err := c.anthropic.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.config.Anthropic.AIModel,
		MaxTokens: maxTokens,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: requestMessages,
	})
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	parts := make([]string, 0, len(resp.Content))
	for _, block := range resp.Content {
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok && textBlock.Text != "" {
			parts = append(parts, textBlock.Text)
		}
	}

	answers := []string{}
	if len(parts) > 0 {
		answers = append(answers, strings.Join(parts, ""))
	}

	return answers, nil
}
