package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"bot/internal/client/github"
	"bot/internal/model"
)

func (s *Service) ProcessPrompt(ctx context.Context, prompt string, threadMessages []*model.ThreadMessage) ([]*Answer, error) {
	resp, err := s.chatCompletion(ctx, prompt, threadMessages)
	if err != nil {
		s.logger.Errorf("chat completion: %s", err.Error())
		return nil, fmt.Errorf("chatCompletion: %w", err)
	}

	if resp.Error != nil {
		s.logger.Errorf("chat completion error: %s (type: %s, param: %s, code: %s)",
			resp.Error.Message, resp.Error.Type, resp.Error.Param, resp.Error.Code)
		return []*Answer{{Answer: resp.Error.Message}}, nil
	}

	result := []*Answer{}
	for _, choice := range resp.Choices {
		if choice.Message == nil {
			continue
		}

		for _, toolCall := range choice.Message.ToolCalls {
			answer, err := s.executeFunction(toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				s.logger.Errorf("execute function %s: %s", toolCall.Function.Name, err.Error())
				continue
			}
			if answer != nil {
				result = append(result, answer)
			}
		}
	}

	return result, nil
}

func (s *Service) chatCompletion(ctx context.Context, prompt string, threadMessages []*model.ThreadMessage) (*github.ChatCompletionResponse, error) {
	messages := make([]*github.ChatCompletionRequestMessage, 0, len(threadMessages)+2)
	messages = append(messages, &github.ChatCompletionRequestMessage{
		Role:    RoleSystem,
		Content: s.config.SystemPrompt,
	})

	for _, threadMessage := range threadMessages {
		switch threadMessage.Type {
		case model.MessageTypeRequest:
			messages = append(messages, &github.ChatCompletionRequestMessage{
				Role:    RoleUser,
				Content: threadMessage.Text,
			})
		case model.MessageTypeResponse:
			messages = append(messages, &github.ChatCompletionRequestMessage{
				Role:    RoleAssistant,
				Content: threadMessage.Text,
			})
		}
	}

	messages = append(messages, &github.ChatCompletionRequestMessage{
		Role:    RoleUser,
		Content: prompt,
	})

	var err error
	for _, AIModel := range s.config.AIModels {
		var resp *github.ChatCompletionResponse
		resp, err = s.githubClient.ChatCompletions(ctx, &github.ChatCompletionRequest{
			Model:       AIModel,
			Messages:    messages,
			Tools:       s.tools,
			ToolChoice:  ToolChoiceAuto,
			Temperature: Temperature,
		})
		if err != nil {
			s.logger.Errorf("get chat completion for model %s: %s", AIModel, err.Error())
			continue
		}

		return resp, nil
	}

	return nil, err
}

func (s *Service) executeFunction(functionName, arguments string) (*Answer, error) {
	switch functionName {
	case "send_answer":
		answer := &Answer{}
		if err := json.Unmarshal([]byte(arguments), answer); err != nil {
			s.logger.Errorf("unmarshal answer: %s", err.Error())
			return nil, err
		}
		return answer, nil
	default:
		s.logger.Warnf("unknown function name: %s", functionName)
	}

	return nil, nil
}
