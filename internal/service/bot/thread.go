package bot

import (
	"bot/internal/model"
	"context"
	"time"
)

func (s *Service) getThread(threadID int64) *model.Thread {
	thread, exists := s.threads[threadID]
	if !exists {
		thread = &model.Thread{
			StartedAt: time.Now().UTC(),
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		s.threads[threadID] = thread
	}

	return thread
}

func (s *Service) addThreadMessage(threadID int64, messageType model.MessageType, messageText string) {
	thread := s.getThread(threadID)

	s.mu.Lock()
	defer s.mu.Unlock()

	thread.Messages = append(thread.Messages, &model.ThreadMessage{
		Type: messageType,
		Text: messageText,
	})
}

func (s *Service) getThreadMessages(threadID int64) []*model.ThreadMessage {
	thread, exists := s.threads[threadID]
	if !exists {
		return []*model.ThreadMessage{}
	}

	return thread.Messages
}

func (s *Service) clearThreads(ctx context.Context) {
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

func (s *Service) clearThreadsInt() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for threadID, thread := range s.threads {
		if time.Since(thread.StartedAt) > 24*time.Hour {
			delete(s.threads, threadID)
		}
	}
}
