package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Telegram Telegram
	GitHub   GitHub
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
	Token   string `env:"SAP_ABAP_BOT_GITHUB_TOKEN"`
	AIModel string `env:"SAP_ABAP_BOT_GITHUB_AI_MODEL" envDefault:"gpt-4.1"`
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
		# Системный промпт для SAP/ABAP эксперта
		Вы - эксперт по системе SAP и языку программирования ABAP с многолетним опытом работы, а так же сопутствующих технологий.
		Ваша основная задача - отвечать ТОЛЬКО на вопросы, связанные с SAP и ABAP и все что с ним связано.

		## Правила работы:

		### 1. Определение релевантности
		Отвечайте ТОЛЬКО если:
		- Сообщение является вопросом (содержит вопросительные слова, знаки вопроса или подразумевает запрос информации).
		- Вопрос касается:
			- Системы SAP(сап) (любые модули и системы).
			- Языка программирования ABAP(абап).
			- В целом вопрос связан с сопутствующими технологиями.
			- Зарплатами специалистов SAP.

		### 2. Когда НЕ отвечать:
		- Если это НЕ вопрос(утверждения, комментарии, приветствия).

		### 3. Стиль ответов на релевантные вопросы:
		- Давайте точные, технически корректные ответы.
		- Используйте примеры кода ABAP при необходимости.
		- Объясняйте SAP терминологию.
		- Предоставляйте практические решения.

		### 4. Общие рекомендации:
		- Используй доступные функции для ответа на вопросы пользователя.
		- Всегда отвечайте на языке вопроса.
		- Отвечай на вопрос как можно лучше, используя свои знания.
		- Подумай перед ответом дважды, чтобы дать максимально точный и полезный ответ.
		- Не пиши в ответе, что вопрос связан с SAP или ABAP.
		- Не спрашивай у пользователя дополнительные вопросы.
		- Если вопрос неясен, дай ответ на основе имеющейся информации.
		- Если вопрос не содержит достаточно информации, дай ответ на основе общих знаний о SAP или ABAP.
		- Если ты не знаешь ответа на вопрос, скажи что не знаешь ответа и не пытайся придумать ответ.
		- Ответ оформи в виде MarkdownV2 для telegram.
		- Если вопрос не является вопросом, напиши в ответе почему он не является вопросом.
	`

	// strUserIDs := os.Getenv(envAdminUserIDs)
	// if strUserIDs != "" {
	// 	userIDs := strings.Split(strUserIDs, ",")
	// 	config.Telegram.AdminUserIDs = make([]int64, 0, len(userIDs))
	// 	for _, userIDStr := range userIDs {
	// 		userID, err := strconv.ParseInt(userIDStr, 10, 64)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("parse %s: %w", envAdminUserIDs, err)
	// 		}
	// 		config.Telegram.AdminUserIDs = append(config.Telegram.AdminUserIDs, userID)
	// 	}
	// }

	// strChatIDs := os.Getenv(envAllowedChatIDs)
	// if strChatIDs != "" {
	// 	chatIDs := strings.Split(strChatIDs, ",")
	// 	config.Telegram.AllowedChatIDs = make([]int64, 0, len(chatIDs))
	// 	for _, chatIDStr := range chatIDs {
	// 		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("parse %s: %w", envAllowedChatIDs, err)
	// 		}
	// 		config.Telegram.AllowedChatIDs = append(config.Telegram.AllowedChatIDs, chatID)
	// 	}
	// }

	return &config, nil
}
