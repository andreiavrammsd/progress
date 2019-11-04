package progress

import (
	"io"
	"os"
	"sync"
)

type Config struct {
	Drawer Drawer
	Writer io.Writer
	Error  Error
}

type Error func(error)

type Progress struct {
	drawer Drawer
	writer io.Writer
	tick   chan struct{}
	stop   chan struct{}
	error  Error
	lock   sync.Mutex
}

func (p *Progress) Progress() {
	p.lock.Lock()
	if p.tick != nil {
		p.tick <- struct{}{}
	}
	p.lock.Unlock()
}

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
