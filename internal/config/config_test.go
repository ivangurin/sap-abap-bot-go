package config

import (
	"strings"
	"testing"
)

// В лог не должны попадать токены целиком.
func TestLogStringMasksTokens(t *testing.T) {
	t.Parallel()

	c := &Config{}
	c.App.AIProvider = "github"
	c.App.LogLevel = "debug"
	c.Telegram.BotToken = "1234567890:AA-telegram-secret-token"
	c.Telegram.AdminUserIDs = []int64{1, 2}
	c.Telegram.AllowedChatIDs = []int64{-100}
	c.GitHub.Token = "github_pat_super-secret-value"
	c.Anthropic.Token = "sk-ant-api03-super-secret"

	out := c.LogString()

	for _, secret := range []string{
		c.Telegram.BotToken,
		c.GitHub.Token,
		c.Anthropic.Token,
	} {
		if strings.Contains(out, secret) {
			t.Fatalf("LogString() содержит токен целиком: %q\n%s", secret, out)
		}
	}

	// Премисс: замаскированные хвосты всё же видны, чтобы отличить заданный токен от пустого.
	for _, masked := range []string{"123…ken", "git…lue", "sk-…ret"} {
		if !strings.Contains(out, masked) {
			t.Fatalf("LogString() не содержит маску %q:\n%s", masked, out)
		}
	}
}

func TestMaskToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		token string
		want  string
	}{
		{"", ""},
		{"short", "***"},
		{"1234567890abcdef", "123…def"},
	}

	for _, c := range cases {
		if got := maskToken(c.token); got != c.want {
			t.Fatalf("maskToken(%q) = %q, хочу %q", c.token, got, c.want)
		}
	}
}
