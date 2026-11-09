package bot

import (
	"sync"
	"testing"

	tgmodels "github.com/go-telegram/bot/models"

	"bot/internal/model"
)

func newTestService() *service {
	return &service{
		threads: make(map[threadKey]*model.Thread),
	}
}

// getThreadMessages должен возвращать копию: последующие добавления в тред
// не должны менять уже полученный срез (иначе гонка слайсов при итерации).
func TestGetThreadMessagesReturnsCopy(t *testing.T) {
	t.Parallel()

	s := newTestService()
	key := threadKey{chatID: 10, threadID: 20}

	s.addThreadMessage(key, model.MessageTypeRequest, "вопрос")

	before := s.getThreadMessages(key)
	s.addThreadMessage(key, model.MessageTypeResponse, "ответ")

	if len(before) != 1 {
		t.Fatalf("getThreadMessages вернул живой слайс: len = %d, хочу 1 (копия на момент вызова)", len(before))
	}
	if after := s.getThreadMessages(key); len(after) != 2 {
		t.Fatalf("getThreadMessages после добавления: len = %d, хочу 2", len(after))
	}
}

// Личный чат: один тред на чат — диалог в личке не теряется между сообщениями.
func TestThreadKeyPrivateChat(t *testing.T) {
	t.Parallel()

	msg := &tgmodels.Message{
		ID: 100,
		Chat: tgmodels.Chat{
			ID:   42,
			Type: tgmodels.ChatTypePrivate,
		},
		From: &tgmodels.User{ID: 42},
	}

	key := threadKeyFromMessage(msg)

	if key != (threadKey{chatID: 42, threadID: 0}) {
		t.Fatalf("threadKeyFromMessage() = %+v, хочу {chatID:42 threadID:0}", key)
	}
}

// Форум-топик: тред привязан к топику внутри чата.
func TestThreadKeyForumTopic(t *testing.T) {
	t.Parallel()

	msg := &tgmodels.Message{
		MessageThreadID: 7,
		Chat: tgmodels.Chat{
			ID:   -100,
			Type: tgmodels.ChatTypeSupergroup,
		},
		From: &tgmodels.User{ID: 42},
	}

	key := threadKeyFromMessage(msg)

	if key != (threadKey{chatID: -100, threadID: 7}) {
		t.Fatalf("threadKeyFromMessage() = %+v, хочу {chatID:-100 threadID:7}", key)
	}
}

// Обычная группа без топика: тред на конкретное сообщение, но в пределах чата.
func TestThreadKeyGroupMessage(t *testing.T) {
	t.Parallel()

	msg := &tgmodels.Message{
		ID: 100,
		Chat: tgmodels.Chat{
			ID:   -100,
			Type: tgmodels.ChatTypeGroup,
		},
		From: &tgmodels.User{ID: 42},
	}

	key := threadKeyFromMessage(msg)

	if key != (threadKey{chatID: -100, threadID: 100}) {
		t.Fatalf("threadKeyFromMessage() = %+v, хочу {chatID:-100 threadID:100}", key)
	}
}

// Одинаковые ID сообщений/топиков в разных чатах — разные треды.
func TestThreadKeyDifferentChats(t *testing.T) {
	t.Parallel()

	first := &tgmodels.Message{
		ID:   123,
		Chat: tgmodels.Chat{ID: -100},
		From: &tgmodels.User{ID: 42},
	}
	second := &tgmodels.Message{
		ID:   123,
		Chat: tgmodels.Chat{ID: -200},
		From: &tgmodels.User{ID: 42},
	}

	if threadKeyFromMessage(first) == threadKeyFromMessage(second) {
		t.Fatal("сообщения с одинаковым ID из разных чатов дали одинаковый ключ треда")
	}
}

// Конкурентный доступ к хранилищу тредов (go test -race):
// записи, чтения и проверки существования треда из параллельных горутин.
func TestThreadsConcurrentAccess(t *testing.T) {
	t.Parallel()

	s := newTestService()
	key := threadKey{chatID: 10, threadID: 20}

	var wg sync.WaitGroup

	for range 50 {
		wg.Go(func() {
			s.addThreadMessage(key, model.MessageTypeRequest, "вопрос")
		})

		wg.Go(func() {
			_ = s.getThreadMessages(key)
		})

		wg.Go(func() {
			_ = s.hasThread(key)
		})
	}

	wg.Wait()

	if got := len(s.getThreadMessages(key)); got != 50 {
		t.Fatalf("после 50 конкурентных добавлений в треде %d сообщений, хочу 50", got)
	}
}
