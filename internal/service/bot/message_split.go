package bot

import "unicode/utf16"

// maxMessageLength — лимит Telegram на длину сообщения (UTF-16 кодовых юнитов).
const maxMessageLength = 4096

// splitMessage режет текст на куски не длиннее limit UTF-16 кодовых юнитов —
// именно в этих единицах Telegram считает длину сообщения (эмодзи занимает два).
func splitMessage(text string, limit int) []string {
	if text == "" {
		return nil
	}

	chunks := make([]string, 0, 1)

	runes := make([]rune, 0, limit)
	size := 0

	for _, r := range text {
		runeLen := utf16.RuneLen(r)

		if size+runeLen > limit && len(runes) > 0 {
			chunks = append(chunks, string(runes))

			runes = runes[:0]
			size = 0
		}

		runes = append(runes, r)
		size += runeLen
	}

	if len(runes) > 0 {
		chunks = append(chunks, string(runes))
	}

	return chunks
}
