package bot

import (
	"bot/internal/model"
	"bot/internal/service/agent"
	"context"
	"sync"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	pkg_config "bot/internal/config"

	pkg_logger "bot/internal/pkg/logger"
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
	threads      map[int64]*model.Thread
	mu           sync.Mutex
}

func NewService(
	ctx context.Context,
	config *pkg_config.Config,
	logger pkg_logger.Logger,
	agentService agent.Service,
) Service {
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

	var err error
	service.bot, err = tgbot.New(config.Telegram.BotToken, opts...)
	if err != nil {
		panic("failed to create bot: " + err.Error())
	}

	service.threads = make(map[int64]*model.Thread)

	return service
}
