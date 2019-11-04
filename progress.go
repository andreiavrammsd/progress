// Package progress offers an indicator that a certain flow is actually working.
package progress

import (
	"io"
	"os"
	"sync"
)

// Config allows setup of progress.
type Config struct {
	Drawer Drawer
	Writer io.Writer
	Error  Error
}

// Error represents the callback type to pass if you want to see all errors from a Drawer
type Error func(error)

// Progress represents a progress instance.
type Progress struct {
	drawer Drawer
	writer io.Writer
	tick   chan struct{}
	stop   chan struct{}
	error  Error
	lock   sync.Mutex
}

// Progress should be called each time the Drawer should update.
func (p *Progress) Progress() {
	p.lock.Lock()
	if p.tick != nil {
		p.tick <- struct{}{}
	}
	p.lock.Unlock()
}

// Stop cancels progress animation.
func (p *Progress) Stop() {
	close(p.stop)
}

func (p *Progress) start() {
	p.tick = make(chan struct{})
	p.stop = make(chan struct{})

	go func() {
		for {
			select {
			case <-p.tick:
				if err := p.drawer.Draw(p.writer); err != nil && p.error != nil {
					p.error(err)
				}
			case <-p.stop:
				p.lock.Lock()
				close(p.tick)
				p.tick = nil
				p.lock.Unlock()
				return
			}
		}
	}()
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

	p := &Progress{
		drawer: c.Drawer,
		writer: c.Writer,
		error:  c.Error,
	}

	p.start()

	return p
}
