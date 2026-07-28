package bot

import (
	"context"

	pkg_config "bot/internal/config"
	pkg_closer "bot/internal/pkg/closer"
	pkg_logger "bot/internal/pkg/logger"
	service_provider "bot/internal/service/provider"
)

type App struct {
	config *pkg_config.Config
}

func NewApp(config *pkg_config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := pkg_logger.NewLogger(a.config.App.LogLevel, a.config.App.LogFile)

	logger.Info(ctx, "app is starting...")
	defer logger.Info(ctx, "app has been finished")

	logger.Debug(ctx, "with config", "config", a.config)

	closer := pkg_closer.New(ctx, logger)
	defer func() {
		err := closer.Close()
		if err != nil {
			logger.Info(ctx, "errors on close", "error", err.Error())
		}
	}()

	sp := service_provider.NewProvider(ctx, a.config, logger)

	errCh := make(chan error, 1)
	go func() {
		logger.Info(ctx, "starting bot service...")
		err := sp.GetBotService().Run(ctx)
		if err != nil {
			logger.Error(ctx, "run bot", "error", err.Error())
			errCh <- err
			closer.Stop()
		}
	}()

	closer.Add(func() error {
		logger.Info(ctx, "stopping bot service ...")
		err := sp.GetBotService().Close(ctx)
		if err != nil {
			logger.Error(ctx, "stop bot service", "error", err.Error())
			return err
		}
		logger.Info(ctx, "bot service has been finished")
		return nil
	})

	closer.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
