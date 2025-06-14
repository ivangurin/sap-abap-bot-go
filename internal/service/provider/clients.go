package service_provider

import "bot/internal/client/github"

type clients struct {
	githubClient github.IClient
}

func (sp *Provider) GetGithubClient() github.IClient {
	if sp.githubClient == nil {
		sp.githubClient = github.NewClient(
			sp.config,
			sp.logger,
		)
	}
	return sp.githubClient
}
