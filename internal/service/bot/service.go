package bot

import (
	"context"
	"sync"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	pkg_config "bot/internal/config"
	"bot/internal/model"
	pkg_logger "bot/internal/pkg/logger"
	"bot/internal/service/agent"
)

type IService interface {
	Run(ctx context.Context) error
	DefaultHandler(ctx context.Context, bot *tgbot.Bot, update *models.Update)
	ErrorHandler(err error)
	Close(ctx context.Context) error
}

type Service struct {
	config       *pkg_config.Config
	logger       *pkg_logger.Logger
	bot          *tgbot.Bot
	agentService agent.IService
	username     string
	threads      map[int64]*model.Thread
	mu           sync.Mutex
}

func NewService(
	config *pkg_config.Config,
	logger *pkg_logger.Logger,
	agentService agent.IService,
) IService {
	service := &Service{
		config:       config,
		logger:       logger,
		agentService: agentService,
	}

	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(service.DefaultHandler),
		tgbot.WithErrorsHandler(service.ErrorHandler),
	}

	if config.Debug {
		opts = append(opts, tgbot.WithDebug())
	}

	var err error
	service.bot, err = tgbot.New(config.BotToken, opts...)
	if err != nil {
		logger.Fatalf("create bot: %s", err.Error())
	}

	service.threads = make(map[int64]*model.Thread)

	return service
}
