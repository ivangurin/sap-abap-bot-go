package bot

import (
	"context"
	"slices"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"bot/internal/model"
)

// Служебные сообщения: если к боту обратились, он обязан ответить в любом случае.
const (
	errNoAnswerText = "Не удалось получить ответ от ИИ, попробуйте ещё раз."
	emptyAnswerText = "ИИ вернул пустой ответ, попробуйте переформулировать вопрос."
)

// parseModeNone — отправка без форматирования.
const parseModeNone = models.ParseMode("")

func (s *service) DefaultHandler(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if !s.shouldHandle(update) {
		return
	}

	msg := update.Message

	thread := threadKeyFromMessage(msg)

	messageText := strings.Replace(msg.Text, "@"+s.username, "", 1)

	threadMessages := s.getThreadMessages(thread)

	answers, err := s.agentService.ProcessPrompt(ctx, messageText, threadMessages)
	if err != nil {
		s.logger.Errorf(ctx, "process prompt: %s", err.Error())

		s.notifyUser(ctx, bot, msg.Chat.ID, msg.ID, errNoAnswerText)

		return
	}

	// Добавляем вопрос в тред
	s.addThreadMessage(thread, model.MessageTypeRequest, messageText)

	if len(answers) == 0 {
		s.notifyUser(ctx, bot, msg.Chat.ID, msg.ID, emptyAnswerText)

		return
	}

	for _, answer := range answers {
		s.sendAnswer(ctx, bot, msg.Chat.ID, msg.ID, answer.Answer)

		// Добавляем ответ в тред
		s.addThreadMessage(thread, model.MessageTypeResponse, answer.Answer)
	}
}

// shouldHandle решает, реагирует ли бот на обновление: на формат и автора
// сообщения, а также на правила гейтинга для лички и групп. При ответе
// на упоминание подтягивает текст сообщения, на которое ответили.
func (s *service) shouldHandle(update *models.Update) bool {
	if update.Message == nil || update.Message.From == nil || update.Message.From.IsBot {
		return false
	}

	msg := update.Message

	// Нетекстовые обновления (стикеры, фото, голосовые) бот не обрабатывает.
	if msg.Text == "" {
		return false
	}

	// Если сообщение в чате с ботом — отвечаем только админам
	if msg.Chat.ID == msg.From.ID {
		return slices.Contains(s.config.Telegram.AdminUserIDs, msg.From.ID)
	}

	// Если сообщение в группе — только из разрешённых чатов
	if !slices.Contains(s.config.Telegram.AllowedChatIDs, msg.Chat.ID) {
		return false
	}

	// Бот уже участвует в переписке этого треда
	if s.hasThread(threadKeyFromMessage(msg)) {
		return true
	}

	// Или его упомянули; текст сообщения, на которое ответили, — часть вопроса
	if strings.Contains(msg.Text, "@"+s.username) {
		if msg.ReplyToMessage != nil {
			msg.Text = msg.ReplyToMessage.Text + "\n" + msg.Text
		}

		return true
	}

	return false
}

// sendAnswer отправляет ответ пользователю: нарезает длинный текст на части
// (лимит Telegram — 4096 UTF-16 юнитов), а при неудачной отправке с
// форматированием повторяет часть без него — ответ обязан дойти в любом виде.
func (s *service) sendAnswer(ctx context.Context, bot *tgbot.Bot, chatID int64, replyToMessageID int, answer string) {
	for _, chunk := range splitMessage(answer, maxMessageLength) {
		err := s.sendMessage(ctx, bot, chatID, replyToMessageID, formatAnswer(chunk), models.ParseModeMarkdown)
		if err != nil {
			s.logger.Warnf(ctx, "send message with markdown: %s", err.Error())

			err = s.sendMessage(ctx, bot, chatID, replyToMessageID, chunk, parseModeNone)
			if err != nil {
				s.logger.Errorf(ctx, "send message fallback: %s", err.Error())
			}
		}
	}
}

// notifyUser отправляет короткое служебное сообщение без форматирования.
func (s *service) notifyUser(ctx context.Context, bot *tgbot.Bot, chatID int64, replyToMessageID int, text string) {
	err := s.sendMessage(ctx, bot, chatID, replyToMessageID, text, parseModeNone)
	if err != nil {
		s.logger.Errorf(ctx, "send notification: %s", err.Error())
	}
}

func (s *service) sendMessage(
	ctx context.Context,
	bot *tgbot.Bot,
	chatID int64,
	replyToMessageID int,
	text string,
	parseMode models.ParseMode,
) error {
	params := &tgbot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMessageID,
		},
	}
	if parseMode != parseModeNone {
		params.ParseMode = parseMode
	}

	_, err := bot.SendMessage(ctx, params)

	return err
}
