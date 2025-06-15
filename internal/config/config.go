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

	config := &Config{
		AIModel:      "gpt-4.1",
		SystemPrompt: "Ты эксперт по системам SAP и программированию на языке ABAP. Определи связан ли вопрос с системами SAP или ABAP. Если вопрос связан, используй доступные функции для ответа на вопросы пользователя. Отвечай на вопрос на том языке, на котором задан вопрос. Отвечай на вопрос как можно лучше, используя свои знания. Если вопрос связан с ABAP, то приведи пример кода. Если вопрос не связан с системами SAP или ABAP, укажи что вопрос не связан с SAP и ABAP и не отвечай на него. Не спрашивай у пользователя дополнительные вопросы, просто отвечай на вопрос. Если ты не знаешь ответа на вопрос, скажи что не знаешь ответа и не пытайся придумать ответ.",
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
