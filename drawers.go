package progress

import (
	"io"
)

const clear = "\033[H\033[2J"

// Drawer is the interface any progress drawer must implement.
type Drawer interface {
	Draw(io.Writer) error
}

// Characters draws any given characters (Default: ".")
type Characters struct {
	Characters []byte
}

// Draw draws Characters.
func (c *Characters) Draw(w io.Writer) (err error) {
	if len(c.Characters) == 0 {
		c.Characters = []byte(".")
	}

	_, err = w.Write(c.Characters)
	return
}

// Spinner draws a spinning progress
type Spinner struct {
	index int
	signs []string
}

// Draw draws Spinner.
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

// Blink draws a blink effect.
type Blink struct {
	Characters []byte
	chars      [][]byte
	index      uint
}

// Draw draws Blink.
func (b *Blink) Draw(w io.Writer) (err error) {
	if len(b.Characters) == 0 {
		b.Characters = []byte("*")
	}

	if b.chars == nil {
		b.chars = [][]byte{
			b.Characters,
			[]byte(clear),
		}
	}

	b.index ^= 1
	_, err = w.Write(b.chars[b.index])

	return
}

// Arrow draws an arrow sign.
type Arrow struct {
	data []byte
}

// Draw draws Arrow.
func (b *Arrow) Draw(w io.Writer) (err error) {
	if b.data == nil {
		b.data = []byte("=>\033[1D")
	}
	_, err = w.Write(b.data)
	return
}
