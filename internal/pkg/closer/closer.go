package closer

import (
	"bot/internal/pkg/logger"
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	errLogKey = "error"
	errMsg    = "close resource error"
)

// Closer собирает функции закрытия ресурсов и выполняет их в обратном порядке (LIFO)
// при вызове Close. Метод Wait блокируется до сигнала завершения (SIGTERM/SIGINT)
// либо до программной отмены через Stop.
type Closer struct {
	ctx      context.Context
	log      logger.Logger
	mu       sync.Mutex
	funcs    []func() error
	closed   bool
	stop     chan struct{}
	stopOnce sync.Once
}

// New создаёт новый Closer с переданным логгером.
func New(ctx context.Context, log logger.Logger) *Closer {
	return &Closer{
		ctx:   ctx,
		log:   log,
		funcs: make([]func() error, 0),
		stop:  make(chan struct{}),
	}
}

// Add добавляет функцию закрытия ресурса.
// Потокобезопасен: может вызываться из разных горутин.
// Если Add вызван после Close, функция f будет выполнена немедленно,
// а в лог будет записано предупреждение (чтобы избежать утечки ресурса).
// Если f == nil, функция игнорируется с предупреждением в лог.
func (c *Closer) Add(f func() error) {
	if f == nil {
		c.log.Warn(c.ctx, "attempted to add nil function to closer, ignoring")

		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		// Если уже закрыто, выполняем немедленно и логируем предупреждение
		c.log.Warn(c.ctx, "add called after close: executing immediately to prevent resource leak")
		_ = f()

		return
	}

	c.funcs = append(c.funcs, f)
}

// Wait блокируется до получения сигнала завершения (SIGTERM или SIGINT) либо до вызова Stop,
// после чего возвращает управление. Сами ресурсы закрывает отложенный Close.
// Сигналы слушаются в канале с буфером 1, чтобы не пропустить сигнал
// при задержке между Notify и чтением из канала.
func (c *Closer) Wait() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sig)

	select {
	case <-sig:
		c.log.Info(c.ctx, "received termination signal")
	case <-c.stop:
		// Программная отмена через Stop — например, при фатальной ошибке старта сервера.
	}
}

// Stop программно разблокирует Wait, не дожидаясь сигнала.
// Идемпотентен. Применяется, когда ждать сигнал не нужно — например,
// при фатальной ошибке запуска критического компонента (http server и т.п.).
func (c *Closer) Stop() {
	c.stopOnce.Do(func() {
		close(c.stop)
	})
}

// Close выполняет все добавленные функции в обратном порядке (LIFO):
// последняя добавленная функция выполняется первой.
// Идемпотентен: повторный вызов не выполняет функции повторно и возвращает nil.
// Мьютекс НЕ удерживается во время вызова функций, чтобы избежать deadlock
// или долгой блокировки при ошибках/зависаниях.
// Каждая ошибка логируется, метод возвращает агрегированную ошибку через errors.Join.
// Также разблокирует возможные ожидающие вызовы Wait.
func (c *Closer) Close() error {
	// На случай, если Close вызывают до возврата Wait — разблокируем ждущих.
	c.stopOnce.Do(func() {
		close(c.stop)
	})

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()

		return nil
	}
	c.closed = true

	// Копируем функции, чтобы выполнять без удержания мьютекса
	funcs := make([]func() error, len(c.funcs))
	copy(funcs, c.funcs)

	c.mu.Unlock()

	// Выполняем в обратном порядке (LIFO)
	var errs []error
	for i := len(funcs) - 1; i >= 0; i-- {
		if err := funcs[i](); err != nil {
			c.log.Error(c.ctx, errMsg, map[string]any{errLogKey: err})
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
