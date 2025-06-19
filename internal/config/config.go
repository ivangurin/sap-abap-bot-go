package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// nolint: gosec
const (
	envBotToken       = "SAP_ABAP_BOT_TOKEN"
	envGitHubToken    = "SAP_ABAP_BOT_GITHUB_TOKEN"
	envAdminUserIDs   = "SAP_ABAP_BOT_ADMIN_USER_IDS"
	envAllowedChatIDs = "SAP_ABAP_BOT_ALLOWED_CHAT_IDS"
	envDebug          = "SAP_ABAP_BOT_DEBUG"
)

type Config struct {
	AIModel        string
	SystemPrompt   string
	BotToken       string
	GitHubToken    string
	AdminUserIDs   []int64
	AllowedChatIDs []int64
	Debug          bool
}

func NewConfig() (*Config, error) {
	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			return nil, fmt.Errorf("load .env file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("check .env file: %w", err)
	}

	systemPrompt :=
		`
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

	config := &Config{
		AIModel:      "gpt-4.1-mini",
		SystemPrompt: systemPrompt,
		BotToken:     os.Getenv(envBotToken),
		GitHubToken:  os.Getenv(envGitHubToken),
	}

	strUserIDs := os.Getenv(envAdminUserIDs)
	if strUserIDs != "" {
		userIDs := strings.Split(strUserIDs, ",")
		config.AdminUserIDs = make([]int64, 0, len(userIDs))
		for _, userIDStr := range userIDs {
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse %s: %w", envAdminUserIDs, err)
			}
			config.AdminUserIDs = append(config.AdminUserIDs, userID)
		}
	}

	strChatIDs := os.Getenv(envAllowedChatIDs)
	if strChatIDs != "" {
		chatIDs := strings.Split(strChatIDs, ",")
		config.AllowedChatIDs = make([]int64, 0, len(chatIDs))
		for _, chatIDStr := range chatIDs {
			chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse %s: %w", envAllowedChatIDs, err)
			}
			config.AllowedChatIDs = append(config.AllowedChatIDs, chatID)
		}
	}

	debug := os.Getenv(envDebug)
	if debug != "" {
		debugValue, err := strconv.ParseBool(debug)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", envDebug, err)
		}
		config.Debug = debugValue
	}

	return config, nil
}
