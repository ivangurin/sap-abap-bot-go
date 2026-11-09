package bot

import (
	"strings"

	tgbot "github.com/go-telegram/bot"
)

// formatAnswer готовит ответ модели к отправке в Telegram с MarkdownV2:
// экранирует спецсимволы и возвращает разрешённое системным промптом
// форматирование — жирный текст и блоки кода в тройных кавычках.
//
// Жирный текст в MarkdownV2 — одиночная звёздочка (*текст*), модель же
// по промпту пишет **текст**, поэтому пара экранированных звёздочек
// заменяется на одну.
func formatAnswer(answer string) string {
	formatted := tgbot.EscapeMarkdownUnescaped(answer)
	formatted = strings.ReplaceAll(formatted, "\\`\\`\\`", "```")
	formatted = strings.ReplaceAll(formatted, "\\*\\*", "*")

	return formatted
}
