package service_provider

import (
	"bot/internal/clients/anthropic"
	"bot/internal/clients/github"
)

type clients struct {
	githubClient    github.Client
	anthropicClient anthropic.Client
}

func (sp *Provider) GetGithubClient() github.Client {
	if sp.githubClient == nil {
		sp.githubClient = github.NewClient(
			sp.config,
			sp.logger,
		)
	}
	return sp.githubClient
}

func (sp *Provider) GetAnthropicClient() anthropic.Client {
	if sp.anthropicClient == nil {
		sp.anthropicClient = anthropic.NewClient(
			sp.config,
			sp.logger,
		)
	}
	return sp.anthropicClient
}
