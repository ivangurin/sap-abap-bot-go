package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"bot/internal/client/github"
)

func (s *Service) ProcessPrompt(ctx context.Context, prompt string) ([]*Answer, error) {
	resp, err := s.chatCompletion(ctx, prompt)
	if err != nil {
		s.logger.Errorf("get chat completion: %s", err.Error())
		return nil, fmt.Errorf("chatCompletion: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("chat completion error: %s (type: %s, param: %s, code: %s)",
			resp.Error.Message, resp.Error.Type, resp.Error.Param, resp.Error.Code)
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

func (s *Service) chatCompletion(ctx context.Context, prompt string) (*github.ChatCompletionResponse, error) {
	var err error
	for _, model := range s.config.AIModels {
		var resp *github.ChatCompletionResponse
		resp, err = s.githubClient.ChatCompletions(ctx, &github.ChatCompletionRequest{
			Model: model,
			Messages: []*github.ChatCompletionRequestMessage{
				{
					Role:    "system",
					Content: s.config.SystemPrompt,
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Tools:       s.tools,
			ToolChoice:  "auto",
			Temperature: 0.7,
		})
		if err != nil {
			s.logger.Errorf("get chat completion for model %s: %s", model, err.Error())
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
