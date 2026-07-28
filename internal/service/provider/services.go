package service_provider

import (
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
			sp.GetGithubClient(),
		)
	}
	return sp.agentService
}
