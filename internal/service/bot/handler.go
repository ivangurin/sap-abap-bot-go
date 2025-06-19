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
		if !slices.Contains(s.config.AdminUserIDs, update.Message.From.ID) {
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
		s.logger.Infof("prompt: %s, answer: %s (answered: %t, chatID: %d)", update.Message.Text, answer.Answer, answer.Answered, update.Message.Chat.ID)
		if !answer.Answered {
			continue
		}
		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID:          update.Message.Chat.ID,
			MessageThreadID: update.Message.MessageThreadID,
			Text:            tgbot.EscapeMarkdownUnescaped(answer.Answer),
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
