package closer

import (
	pkg_logger "bot/internal/pkg/logger"
	"os"
	"os/signal"
	"sync"
)

type Closer interface {
	Add(f ...func() error)
	Wait()
	CloseAll()
}

type closer struct {
	logger *pkg_logger.Logger
	sync.Mutex
	once     sync.Once
	done     chan struct{}
	funcs    []func() error
	shutdown chan os.Signal
}

// os.Interrupt, syscall.SIGINT, syscall.SIGTERM
func NewCloser(logger *pkg_logger.Logger, sig ...os.Signal) Closer {
	closer := &closer{
		logger:   logger,
		done:     make(chan struct{}),
		shutdown: make(chan os.Signal, 1),
	}

	if len(sig) > 0 {
		go func() {
			signal.Notify(closer.shutdown, sig...)
			<-closer.shutdown
			signal.Stop(closer.shutdown)
			closer.logger.Info("graceful shutdown started...")
			defer closer.logger.Info("graceful shutdown finished")
			closer.CloseAll()
		}()
	}

	return closer
}

func (c *closer) Add(f ...func() error) {
	c.Lock()
	c.funcs = append(c.funcs, f...)
	c.Unlock()
}

func (c *closer) Wait() {
	// if r := recover(); r != nil {
	// 	c.logger.Errorf("panic while waiting graceful shutdown: %v", r)
	// 	c.CloseAll()
	// 	return
	// }
	<-c.done
}

func (c *closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.Lock()
		funcs := c.funcs
		c.Unlock()

		for i := len(funcs) - 1; i >= 0; i-- {
			err := c.funcs[i]()
			if err != nil {
				c.logger.Errorf("close some func from shutdown: %s", err.Error())
			}
		}
	})
}
