package bot

import (
	"context"
	"slices"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Service) Handler(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.From.IsBot {
		return
	}
	if update.Message.ReplyToMessage != nil {
		return
	}

	if update.Message.Chat.ID == update.Message.From.ID {
		if update.Message.From.ID != *s.config.AdminUserID {
			return
		}
	} else {
		if s.config.AllowedChatIDs != nil {
			if !slices.Contains(s.config.AllowedChatIDs, update.Message.Chat.ID) {
				return
			}
		}
	}

	answers, err := s.agentService.ProcessPrompt(ctx, update.Message.Text)
	if err != nil {
		s.logger.Errorf("process prompt: %s", err.Error())
		return
	}

	if len(answers) == 0 {
		return
	}

	for _, answer := range answers {
		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID:          update.Message.Chat.ID,
			MessageThreadID: update.Message.MessageThreadID,
			Text:            tgbot.EscapeMarkdownUnescaped(answer),
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			s.logger.Errorf("send message: %s", err.Error())
		}
	}
}
