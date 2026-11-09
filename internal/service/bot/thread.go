package bot

import (
	"context"
	"time"

	"github.com/go-telegram/bot/models"

	"bot/internal/model"
)

// threadKey идентифицирует тред: чат + топик (или конкретное сообщение в чате).
// Ключ включает чат, чтобы одинаковые ID сообщений и топиков из разных чатов
// не пересекались.
type threadKey struct {
	chatID   int64
	threadID int64
}

// threadKeyFromMessage вычисляет ключ треда для сообщения:
//   - форум-топик — тред на топик внутри чата;
//   - личный чат — один тред на чат, чтобы диалог помнил историю;
//   - обычная группа без топика — тред на конкретное сообщение.
func threadKeyFromMessage(msg *models.Message) threadKey {
	if msg.MessageThreadID != 0 {
		return threadKey{chatID: msg.Chat.ID, threadID: int64(msg.MessageThreadID)}
	}

	if msg.Chat.ID == msg.From.ID {
		return threadKey{chatID: msg.Chat.ID, threadID: 0}
	}

	return threadKey{chatID: msg.Chat.ID, threadID: int64(msg.ID)}
}

func (s *service) getThread(key threadKey) *model.Thread {
	s.mu.Lock()
	defer s.mu.Unlock()

	thread, exists := s.threads[key]
	if !exists {
		thread = &model.Thread{
			StartedAt: time.Now().UTC(),
		}

		s.threads[key] = thread
	}

	return thread
}

func (s *service) addThreadMessage(key threadKey, messageType model.MessageType, messageText string) {
	if messageText == "" {
		return
	}

	thread := s.getThread(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	thread.Messages = append(thread.Messages, &model.ThreadMessage{
		Type: messageType,
		Text: messageText,
	})
}

// getThreadMessages возвращает копию сообщений треда: хендлер итерируется
// по ней, пока другие горутины могут дополнять тред.
func (s *service) getThreadMessages(key threadKey) []*model.ThreadMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	thread, exists := s.threads[key]
	if !exists {
		return []*model.ThreadMessage{}
	}

	messages := make([]*model.ThreadMessage, len(thread.Messages))
	copy(messages, thread.Messages)

	return messages
}

func (s *service) hasThread(key threadKey) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.threads[key]

	return exists
}

func (s *service) clearThreads(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.clearThreadsInt()
		}
	}
}

func (s *service) clearThreadsInt() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, thread := range s.threads {
		if time.Since(thread.StartedAt) > 24*time.Hour {
			delete(s.threads, key)
		}
	}
}
