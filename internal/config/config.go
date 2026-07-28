package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

const (
	AIProviderGitHub    = "github"
	AIProviderAnthropic = "anthropic"
)

type Config struct {
	App       App
	Telegram  Telegram
	GitHub    GitHub
	Anthropic Anthropic
}

type App struct {
	AIProvider   string `env:"SAP_ABAP_BOT_AI_PROVIDER" envDefault:"github"`
	SystemPrompt string
	LogLevel     string `env:"SAP_ABAP_BOT_LOG_LEVEL" envDefault:"debug"`
	LogFile      string `env:"SAP_ABAP_BOT_LOG_FILE" envDefault:""`
}

type Telegram struct {
	BotToken       string  `env:"SAP_ABAP_BOT_TOKEN"`
	AdminUserIDs   []int64 `env:"SAP_ABAP_BOT_ADMIN_USER_IDS"`
	AllowedChatIDs []int64 `env:"SAP_ABAP_BOT_ALLOWED_CHAT_IDS"`
}
type GitHub struct {
	BaseURL string `env:"SAP_ABAP_BOT_GITHUB_BASE_URL" envDefault:"https://models.inference.ai.azure.com"`
	Token   string `env:"SAP_ABAP_BOT_GITHUB_TOKEN"`
	AIModel string `env:"SAP_ABAP_BOT_GITHUB_AI_MODEL" envDefault:"gpt-4.1"`
}

type Anthropic struct {
	BaseURL string `env:"SAP_ABAP_BOT_ANTHROPIC_BASE_URL" envDefault:""`
	Token   string `env:"SAP_ABAP_BOT_ANTHROPIC_TOKEN"`
	AIModel string `env:"SAP_ABAP_BOT_ANTHROPIC_AI_MODEL" envDefault:"claude-opus-4-8"`
}

func NewConfig() (*Config, error) {
	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			return nil, fmt.Errorf("load .env file: %w", err)
		}
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	config.App.SystemPrompt = `
		# Роль
		Ты — эксперт по системе SAP и языку программирования ABAP с многолетним опытом, а также по сопутствующим технологиям.

		# Как отвечать
		- Всегда отвечай на сообщение пользователя — к тебе обратились, поэтому ответ обязателен.
		- Давай точные и технически корректные ответы; при необходимости приводи примеры кода ABAP и поясняй терминологию SAP.
		- Отвечай на языке сообщения.
		- Отвечай кратко и по существу — сразу суть, без вступлений и без упоминаний о том, что вопрос связан с SAP или ABAP.
		- Не задавай пользователю уточняющих вопросов. Если данных не хватает — отвечай на основе общих знаний.
		- Если не знаешь ответа — честно скажи об этом и ничего не выдумывай.
		- Если сообщение не относится к SAP или ABAP, всё равно ответь коротко и по существу.

		# Форматирование
		Ответ отправляется в Telegram. Из форматирования используй только жирный текст в виде **текст** и блоки кода в тройных обратных кавычках с указанием языка. Ссылки, курсив, заголовки и таблицы не используй.
	`

	return &config, nil
}
