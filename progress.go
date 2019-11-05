// Package progress offers a method to display a progress indicating that certain tasks are really running,
// by drawing on the screen various animations. Each step of task should call the Progress method to update
// the animation.
// Tested with `xterm`.
package progress

import (
	"context"
	"io"
	"os"
	"sync"
)

// Config allows setup of progress.
type Config struct {
	// Drawer is the Drawer implementation.
	Drawer Drawer

	// Writer is the writer where Drawer writes content.
	Writer io.Writer

	// Error is the callback to receive Drawer errors.
	Error Error

	// Context is used to stop the progress.
	Context context.Context
}

// Error represents the callback type to pass if you want to see all errors from a Drawer.
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
