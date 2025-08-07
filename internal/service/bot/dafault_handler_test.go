package bot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEscapeMarkdownInlineCode(t *testing.T) {
	t.Parallel()
	input := "before _text_ `code_with_underscore` after."
	expected := "before \\_text\\_ `code_with_underscore` after\\."
	require.Equal(t, expected, escapeMarkdown(input))
}

func TestEscapeMarkdownFencedCode(t *testing.T) {
	t.Parallel()
	input := "before _text_\n```\nfenced_code_with_underscore\n```\nafter."
	expected := "before \\_text\\_\n```\nfenced_code_with_underscore\n```\nafter\\."
	require.Equal(t, expected, escapeMarkdown(input))
}

func TestEscapeMarkdownInlineAndFencedCode(t *testing.T) {
	t.Parallel()
	input := "escape _markdown_ `inline_code` and\n```\nfenced_code_with_underscore\n```\nfinish."
	expected := "escape \\_markdown\\_ `inline_code` and\n```\nfenced_code_with_underscore\n```\nfinish\\."
	require.Equal(t, expected, escapeMarkdown(input))
}
