package service_provider

import (
	"context"

	pkg_config "bot/internal/config"
	pkg_logger "bot/internal/pkg/logger"
)

type Provider struct {
	ctx    context.Context
	config *pkg_config.Config
	logger pkg_logger.Logger
	clients
	services
}

func NewProvider(
	ctx context.Context,
	config *pkg_config.Config,
	logger pkg_logger.Logger,
) *Provider {
	return &Provider{
		ctx:    ctx,
		config: config,
		logger: logger,
	}
}
