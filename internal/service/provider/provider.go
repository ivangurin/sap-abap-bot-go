package service_provider

import (
	pkg_config "bot/internal/config"
	pkg_logger "bot/internal/pkg/logger"
)

type Provider struct {
	config *pkg_config.Config
	logger *pkg_logger.Logger
	clients
	services
}

func NewProvider(
	config *pkg_config.Config,
	logger *pkg_logger.Logger,
) *Provider {
	return &Provider{
		config: config,
		logger: logger,
	}
}
