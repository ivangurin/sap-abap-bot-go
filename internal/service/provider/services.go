package service_provider

import (
	"bot/internal/service/agent"
	"bot/internal/service/bot"
)

type services struct {
	botService   bot.IService
	agentService agent.IService
}

func (sp *Provider) GetBotService() bot.IService {
	if sp.botService == nil {
		sp.botService = bot.NewService(
			sp.config,
			sp.logger,
			sp.GetAgentService(),
		)
	}
	return sp.botService
}

func (sp *Provider) GetAgentService() agent.IService {
	if sp.agentService == nil {
		sp.agentService = agent.NewService(
			sp.config,
			sp.logger,
			sp.GetGithubClient(),
		)
	}
	return sp.agentService
}
