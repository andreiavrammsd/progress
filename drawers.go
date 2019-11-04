package progress

import (
	"io"
)

const clear = "\033[H\033[2J"

type Drawer interface {
	Draw(io.Writer) error
}

type Characters struct {
	Characters []byte
}

func (c *Characters) Draw(w io.Writer) (err error) {
	if len(c.Characters) == 0 {
		c.Characters = []byte(".")
	}

	_, err = w.Write(c.Characters)
	return
}

type Spinner struct {
	index int
	signs []string
}

func (s *Spinner) Draw(w io.Writer) (err error) {
	if s.signs == nil {
		s.signs = []string{"|", "/", "—", "\\"}
	}

	_, err = w.Write([]byte(clear + s.signs[s.index] + "\n"))

	s.index++
	if s.index == len(s.signs) {
		s.index = 0
	}

	return
}

type Arrow struct {
	data []byte
}

func (b *Arrow) Draw(w io.Writer) (err error) {
	if b.data == nil {
		b.data = []byte("=>\033[1D")
	}
	_, err = w.Write(b.data)
	return
}

type Blink struct {
	chars [][]byte
	index uint
}

func (b *Blink) Draw(w io.Writer) (err error) {
	if b.chars == nil {
		b.chars = [][]byte{
			[]byte("*"),
			[]byte(clear),
		}
	}

	b.index ^= 1
	_, err = w.Write(b.chars[b.index])

	return
}
