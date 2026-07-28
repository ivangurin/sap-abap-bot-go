package service_provider

import "bot/internal/clients/github"

type clients struct {
	githubClient github.Client
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
