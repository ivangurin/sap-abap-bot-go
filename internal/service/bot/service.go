package bot

import (
	"context"
	"fmt"
	"sync"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	pkg_config "bot/internal/config"
	"bot/internal/model"
	pkg_logger "bot/internal/pkg/logger"
	"bot/internal/service/agent"
)

type Service interface {
	Run(ctx context.Context) error
	DefaultHandler(ctx context.Context, bot *tgbot.Bot, update *models.Update)
	ErrorHandler(err error)
	Close(ctx context.Context) error
}

type service struct {
	ctx          context.Context
	config       *pkg_config.Config
	logger       pkg_logger.Logger
	bot          *tgbot.Bot
	agentService agent.Service
	username     string
	threads      map[threadKey]*model.Thread
	mu           sync.Mutex
}

func NewService(
	ctx context.Context,
	config *pkg_config.Config,
	logger pkg_logger.Logger,
	agentService agent.Service,
) (Service, error) {
	service := &service{
		ctx:          ctx,
		config:       config,
		logger:       logger,
		agentService: agentService,
	}

	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(service.DefaultHandler),
		tgbot.WithErrorsHandler(service.ErrorHandler),
	}

	if config.App.LogLevel == pkg_logger.LevelDebug {
		opts = append(opts, tgbot.WithDebug())
	}

	bot, err := tgbot.New(config.Telegram.BotToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	service.bot = bot
	service.threads = make(map[threadKey]*model.Thread)

	return service, nil
}
