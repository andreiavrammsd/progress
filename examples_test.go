package progress

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func Example_basic() {
	c := &Config{
		Drawer: &Characters{Characters: []byte(".")},
	}
	p := New(c)

	p.Progress()

	// Output:
	// .
}

func ExampleWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	c := &Config{
		Drawer:  &Characters{Characters: []byte(".")},
		Context: ctx,
	}
	p := New(c)

	p.Progress()
	time.Sleep(time.Millisecond * 500)
	p.Progress()
	p.Progress()

	// Output:
	// .
}

func ExampleWithCancel() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	c := &Config{
		Drawer:  &Characters{Characters: []byte(".")},
		Context: ctx,
	}
	p := New(c)

	p.Progress()
	cancel()

	// Output:
	// .
}

type customDrawer struct {
	chars  []rune
	prefix string
	suffix string
	index  int
}

func (c *customDrawer) Draw(w io.Writer) (err error) {
	_, err = w.Write([]byte(c.prefix + string(c.chars[c.index]) + c.suffix))

	c.index++
	if c.index == len(c.chars) {
		c.index = 0
	}

	if c.index == 1 {
		err = errors.New("err")
	}

	return
}

func Example_custom_drawer() {
	c := Config{
		Drawer: &customDrawer{
			chars:  []rune{'A', 'B', 'C'},
			prefix: "=> ",
			suffix: "\n",
		},
		Writer: os.Stdout,
		Error: func(err error) {
			fmt.Println(err)
		},
	}
	p := New(&c)

	// A running task
	for i := 1; i <= 4; i++ {
		p.Progress()
		time.Sleep(time.Millisecond)
	}

	// Output:
	// => A
	// err
	// => B
	// => C
	// => A
	// err
}
