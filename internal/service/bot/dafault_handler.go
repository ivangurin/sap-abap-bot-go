package bot

import (
	"context"
	"regexp"
	"slices"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"bot/internal/model"
)

func (s *Service) DefaultHandler(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.From.IsBot {
		return
	}

	messageThreadID := int64(update.Message.MessageThreadID)
	if messageThreadID == 0 {
		messageThreadID = int64(update.Message.ID)
	}

	// Если сообщение в чате с ботом
	if update.Message.Chat.ID == update.Message.From.ID {
		if !slices.Contains(s.config.AdminUserIDs, update.Message.From.ID) {
			return
		}
		// Если сообщение в группе
	} else {
		allowed := false
		_, exists := s.threads[messageThreadID]
		if exists {
			allowed = true
		}
		if !allowed {
			if strings.Contains(update.Message.Text, "@"+s.username) {
				allowed = true
				if update.Message.ReplyToMessage != nil {
					update.Message.Text = update.Message.ReplyToMessage.Text + "\n" + update.Message.Text
				}
			}
		}
		if !allowed {
			return
		}
	}

	messageText := strings.Replace(update.Message.Text, "@"+s.username, "", 1)

	threadMessages := s.getThreadMessages(messageThreadID)

	answers, err := s.agentService.ProcessPrompt(ctx, messageText, threadMessages)
	if err != nil {
		s.logger.Errorf("process prompt: %s", err.Error())
		return
	}

	// Добавляем вопроса в тред
	s.addThreadMessage(messageThreadID, model.MessageTypeRequest, messageText)

	for _, answer := range answers {
		answerEsc := escapeMarkdown(answer.Answer)
    
		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   answerEsc,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			s.logger.Errorf("send message: %s", err.Error())
			continue
		}

		// Добавляем ответ в тред
		s.addThreadMessage(messageThreadID, model.MessageTypeResponse, answer.Answer)
	}
}

var markdownCodeRegex = regexp.MustCompile("(?s)```.*?```|`[^`]*`")

func escapeMarkdown(text string) string {
	matches := markdownCodeRegex.FindAllStringIndex(text, -1)
	var result strings.Builder
	last := 0
	for _, m := range matches {
		result.WriteString(tgbot.EscapeMarkdownUnescaped(text[last:m[0]]))
		result.WriteString(text[m[0]:m[1]])
		last = m[1]
	}
	result.WriteString(tgbot.EscapeMarkdownUnescaped(text[last:]))
	return result.String()
}
