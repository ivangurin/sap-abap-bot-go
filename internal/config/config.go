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
	envAdminUserID    = "SAP_ABAP_BOT_ADMIN_USER_ID"
	envAllowedChatIDs = "SAP_ABAP_BOT_ALLOWED_CHAT_IDS"
	envDebug          = "SAP_ABAP_BOT_DEBUG"
)

type Config struct {
	AIModel        string
	SystemPrompt   string
	BotToken       string
	GitHubToken    string
	AdminUserID    *int64
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
		Вы - эксперт по системе SAP и языку программирования ABAP с многолетним опытом работы. 
		Ваша основная задача - отвечать ТОЛЬКО на вопросы, связанные с SAP и ABAP.

		## Правила работы:

		### 1. Определение релевантности
		Отвечайте ТОЛЬКО если:
		- Сообщение является вопросом (содержит вопросительные слова, знаки вопроса или подразумевает запрос информации)
		- Вопрос касается:
		- Системы SAP (любые модули и системы).
		- Языка программирования ABAP.
		- SAP технологий (SAP HANA, Fiori, UI5, BTP, etc.).
		- SAP интеграций и интерфейсов.
		- SAP архитектуры и администрирования.
		- SAP бизнес-процессов.

		### 2. Когда НЕ отвечать:
		- Если это НЕ вопрос (утверждения, комментарии, приветствия).
		- Если вопрос НЕ связан с SAP или ABAP.
		- Если это общие IT-вопросы, не относящиеся к SAP.
		- Если это вопросы о других ERP системах.

		### 3. Области экспертизы:
		- ABAP Objects, ALV, BAPI, RFC, IDoc
		- SAP модули (FI/CO, MM, SD, PP, HR/HCM, etc.)
		- SAP Workflow, Enhancement Framework
		- SAP HANA, CDS Views, AMDP
		- SAP Fiori, UI5, Gateway
		- SAP Cloud Platform (BTP)
		- Customizing и конфигурация SAP

		### 4. Стиль ответов на релевантные вопросы:
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
	`

	config := &Config{
		AIModel:      "gpt-4.1",
		SystemPrompt: systemPrompt,
		BotToken:     os.Getenv(envBotToken),
		GitHubToken:  os.Getenv(envGitHubToken),
	}

	userID := os.Getenv(envAdminUserID)
	if userID != "" {
		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", envAdminUserID, err)
		}
		config.AdminUserID = &id
	}

	strChatIDs := os.Getenv(envAllowedChatIDs)
	if strChatIDs != "" {
		chatIDs := strings.Split(strChatIDs, ",")
		config.AllowedChatIDs = make([]int64, 0, len(chatIDs))
		for _, chatID := range chatIDs {
			parsedID, err := strconv.ParseInt(chatID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse %s: %w", envAllowedChatIDs, err)
			}
			config.AllowedChatIDs = append(config.AllowedChatIDs, parsedID)
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
