package agent

import (
	"bot/internal/model"
	"context"
	"fmt"
)

func (s *service) ProcessPrompt(ctx context.Context, prompt string, threadMessages []*model.ThreadMessage) ([]*Answer, error) {
	messages := make([]*model.ChatMessage, 0, len(threadMessages)+1)

	for _, threadMessage := range threadMessages {
		switch threadMessage.Type {
		case model.MessageTypeRequest:
			messages = append(messages, &model.ChatMessage{
				Role: model.ChatRoleUser,
				Text: threadMessage.Text,
			})
		case model.MessageTypeResponse:
			messages = append(messages, &model.ChatMessage{
				Role: model.ChatRoleAssistant,
				Text: threadMessage.Text,
			})
		}
	}

	messages = append(messages, &model.ChatMessage{
		Role: model.ChatRoleUser,
		Text: prompt,
	})

	answers, err := s.aiClient.Ask(ctx, s.config.App.SystemPrompt, messages)
	if err != nil {
		s.logger.Errorf(ctx, "ask ai provider: %s", err.Error())
		return nil, fmt.Errorf("ask: %w", err)
	}

	result := make([]*Answer, 0, len(answers))
	for _, answer := range answers {
		result = append(result, &Answer{Answer: answer})
	}

	return result, nil
}
