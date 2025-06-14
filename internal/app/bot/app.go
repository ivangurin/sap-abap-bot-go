package bot

import (
	"context"
	"os"
	"syscall"

	pkg_config "bot/internal/config"
	pkg_closer "bot/internal/pkg/closer"
	pkg_logger "bot/internal/pkg/logger"
	service_provider "bot/internal/service/provider"
)

type App struct {
	ctx    context.Context
	config *pkg_config.Config
	logger *pkg_logger.Logger
	closer pkg_closer.Closer
	sp     *service_provider.Provider
}

func NewApp() (*App, error) {
	config, err := pkg_config.NewConfig()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger := pkg_logger.NewLogger(ctx)

	closer := pkg_closer.NewCloser(logger, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	closer.Add(func() error {
		cancel()
		return nil
	})

	sp := service_provider.NewProvider(config, logger)

	return &App{
		ctx:    ctx,
		config: config,
		logger: logger,
		closer: closer,
		sp:     sp,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("app is starting...")
	defer a.logger.Info("app has been finished")

	go func() {
		a.logger.Info("starting bot service...")
		a.sp.GetBotService().Run(a.ctx)
	}()

	a.closer.Add(func() error {
		err := a.sp.GetBotService().Close(a.ctx)
		if err != nil {
			return err
		}
		a.logger.Info("bot service has been finished")
		return nil
	})

	defer a.closer.Wait()

	return nil
}
