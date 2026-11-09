package service_provider

import (
	"fmt"

	pkg_config "bot/internal/config"
	"bot/internal/service/agent"
	"bot/internal/service/bot"
)

type services struct {
	botService   bot.Service
	agentService agent.Service
}

func (sp *Provider) GetBotService() (bot.Service, error) {
	if sp.botService == nil {
		botService, err := bot.NewService(
			sp.ctx,
			sp.config,
			sp.logger,
			sp.GetAgentService(),
		)
		if err != nil {
			return nil, fmt.Errorf("init bot service: %w", err)
		}

		sp.botService = botService
	}

	return sp.botService, nil
}

func (sp *Provider) GetAgentService() agent.Service {
	if sp.agentService == nil {
		sp.agentService = agent.NewService(
			sp.config,
			sp.logger,
			sp.getAIClient(),
		)
	}

	return sp.agentService
}

// getAIClient возвращает клиент ии в зависимости от выбранного провайдера.
func (sp *Provider) getAIClient() agent.Client {
	switch sp.config.App.AIProvider {
	case pkg_config.AIProviderAnthropic:
		return sp.GetAnthropicClient()
	default:
		return sp.GetGithubClient()
	}
}
