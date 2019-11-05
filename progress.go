// Package progress offers an indicator that a certain flow is actually working.
package progress

import (
	"context"
	"io"
	"os"
	"sync"
)

// Config allows setup of progress.
type Config struct {
	Drawer  Drawer
	Writer  io.Writer
	Error   Error
	Context context.Context
}

// Error represents the callback type to pass if you want to see all errors from a Drawer
type Error func(error)

// Progress represents a progress instance.
type Progress struct {
	ch     chan struct{}
	mtx    sync.Mutex
	drawer Drawer
	writer io.Writer
	error  Error
	ctx    context.Context
}

// Progress should be called each time the Drawer should update.
func (p *Progress) Progress() {
	p.mtx.Lock()
	if p.ch != nil {
		p.ch <- struct{}{}
	}
	p.mtx.Unlock()
}

// New creates Progress instance.
func New(c *Config) *Progress {
	if c == nil {
		c = &Config{}
	}

	if c.Drawer == nil {
		c.Drawer = &Spinner{}
	}

	if c.Writer == nil {
		c.Writer = os.Stdout
	}

	if c.Context == nil {
		c.Context = context.Background()
	}

	p := &Progress{
		ch:     make(chan struct{}),
		drawer: c.Drawer,
		writer: c.Writer,
		error:  c.Error,
		ctx:    c.Context,
	}

	go func() {
		for {
			select {
			case <-p.ch:
				if err := p.drawer.Draw(p.writer); err != nil && p.error != nil {
					p.error(err)
				}
			case <-p.ctx.Done():
				p.mtx.Lock()
				close(p.ch)
				p.ch = nil
				p.mtx.Unlock()
				return
			}
		}
	}()

	return p
}
