package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"bot/internal/client/github"
)

func (s *Service) ProcessPrompt(ctx context.Context, prompt string) ([]string, error) {
	resp, err := s.githubClient.ChatCompletions(ctx, &github.ChatCompletionRequest{
		Model: s.config.AIModel,
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
		s.logger.Errorf("get chat completion: %s", err.Error())
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("chat completion error: %s (type: %s, param: %s, code: %s)",
			resp.Error.Message, resp.Error.Type, resp.Error.Param, resp.Error.Code)
	}

	result := []string{}
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
				result = append(result, *answer)
			}
		}
	}

	return result, nil
}

func (s *Service) executeFunction(functionName, arguments string) (*string, error) {
	switch functionName {
	case "send_answer":
		sendAnswer := &SendAnswer{}
		if err := json.Unmarshal([]byte(arguments), sendAnswer); err != nil {
			s.logger.Errorf("unmarshal send_answer arguments: %s", err.Error())
			return nil, err
		}
		if sendAnswer.CorrectQuestion {
			return &sendAnswer.Answer, nil
		}
	default:
		s.logger.Warnf("unknown function name: %s", functionName)
	}

	return nil, nil
}
