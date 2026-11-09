package bot

import (
	"strings"
	"testing"
	"unicode/utf16"
)

func TestSplitMessageShort(t *testing.T) {
	t.Parallel()

	got := splitMessage("короткий ответ", 4096)

	if len(got) != 1 || got[0] != "короткий ответ" {
		t.Fatalf("splitMessage() = %q, хочу один кусок без изменений", got)
	}
}

func TestSplitMessageEmpty(t *testing.T) {
	t.Parallel()

	got := splitMessage("", 4096)

	if len(got) != 0 {
		t.Fatalf("splitMessage(\"\") = %q, хочу пустой список", got)
	}
}

// Длинный текст режется на куски не длиннее лимита (в UTF-16 кодовых юнитах).
func TestSplitMessageLong(t *testing.T) {
	t.Parallel()

	text := strings.Repeat("a", 25)

	got := splitMessage(text, 10)

	if len(got) != 3 {
		t.Fatalf("splitMessage() вернул %d кусков, хочу 3", len(got))
	}
	for i, chunk := range got {
		if len(chunk) > 10 {
			t.Fatalf("кусок %d длиннее лимита: %d", i, len(chunk))
		}
	}
	if got[0]+got[1]+got[2] != text {
		t.Fatal("склейка кусков не равна исходному тексту")
	}
}

// Эмодзи занимает два UTF-16 юнита — лимит должен считаться по ним.
func TestSplitMessageSurrogatePair(t *testing.T) {
	t.Parallel()

	got := splitMessage("ab😀cd", 4)

	if len(got) != 2 {
		t.Fatalf("splitMessage() = %q, хочу 2 куска", got)
	}
	if got[0] != "ab😀" {
		t.Fatalf("первый кусок = %q, хочу %q (4 UTF-16 юнита)", got[0], "ab😀")
	}
	if utf16.RuneLen('😀') != 2 {
		t.Fatal("премисс теста: эмодзи должен занимать 2 юнита")
	}
}
