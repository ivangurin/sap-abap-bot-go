package service_provider

import (
	pkg_config "bot/internal/config"
	"bot/internal/service/agent"
	"bot/internal/service/bot"
)

type services struct {
	botService   bot.Service
	agentService agent.Service
}

func (sp *Provider) GetBotService() bot.Service {
	if sp.botService == nil {
		sp.botService = bot.NewService(
			sp.ctx,
			sp.config,
			sp.logger,
			sp.GetAgentService(),
		)
	}
	return sp.botService
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
