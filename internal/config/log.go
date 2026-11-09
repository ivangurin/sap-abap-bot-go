package config

import "fmt"

// LogString возвращает конфигурацию для логов с замаскированными токенами —
// секреты не должны попадать в логи даже на уровне debug.
func (c *Config) LogString() string {
	return fmt.Sprintf(
		"provider=%s, log_level=%s, log_file=%s, admins=%v, allowed_chats=%v, telegram_token=%s, github_token=%s, anthropic_token=%s",
		c.App.AIProvider,
		c.App.LogLevel,
		c.App.LogFile,
		c.Telegram.AdminUserIDs,
		c.Telegram.AllowedChatIDs,
		maskToken(c.Telegram.BotToken),
		maskToken(c.GitHub.Token),
		maskToken(c.Anthropic.Token),
	)
}

// maskToken маскирует токен, оставляя видными короткие начало и конец,
// чтобы отличить заданный токен от пустого.
func maskToken(token string) string {
	if token == "" {
		return ""
	}

	if len(token) <= 8 {
		return "***"
	}

	return token[:3] + "…" + token[len(token)-3:]
}
