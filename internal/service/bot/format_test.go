package bot

import (
	"strings"
	"testing"
)

// Жирный текст: модель пишет **текст**, в MarkdownV2 жирный — одиночная звёздочка.
func TestFormatAnswerBold(t *testing.T) {
	t.Parallel()

	got := formatAnswer("это **важный** ответ")

	want := "это *важный* ответ"
	if got != want {
		t.Fatalf("formatAnswer() = %q, want %q", got, want)
	}
}

// Блоки кода в тройных кавычках должны доходить до Telegram без экранирования.
func TestFormatAnswerCodeBlock(t *testing.T) {
	t.Parallel()

	got := formatAnswer("вот код:\n```abap\nWRITE: / 'hello'.\n```")

	if !strings.Contains(got, "```abap") {
		t.Fatalf("formatAnswer() = %q, хочу неэкранированное начало блока ```abap", got)
	}
	if !strings.HasSuffix(got, "```") {
		t.Fatalf("formatAnswer() = %q, хочу неэкранированное закрытие блока ```", got)
	}
}

// Спецсимволы вне разрешённого форматирования должны оставаться экранированными.
func TestFormatAnswerEscapesSpecials(t *testing.T) {
	t.Parallel()

	got := formatAnswer("см. раздел 1.2 [тут]!")

	if !strings.Contains(got, `1\.2`) {
		t.Fatalf("formatAnswer() = %q, точка должна быть экранирована", got)
	}
	if !strings.Contains(got, `\[тут\]`) {
		t.Fatalf("formatAnswer() = %q, скобки должны быть экранированы", got)
	}
}
