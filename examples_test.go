package progress

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
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

func Example_advanced() {
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	c := &Config{
		Drawer:  &Characters{},
		Writer:  nil,
		Error:   nil,
		Context: ctx,
	}
	p := New(c)

	task := func(p *Progress) {
		defer close(done)
		wg := sync.WaitGroup{}

		for i := 0; i <= 1; i++ {
			wg.Add(1)

			go func(p *Progress) {
				defer wg.Done()

				counter := 0
				for {
					if counter == 3 {
						break
					}
					counter++

					// Long work
					time.Sleep(time.Millisecond * 50)

					// Progress step
					p.Progress()
				}
			}(p)
		}

		wg.Wait()
	}

	go task(p)
	<-done
	cancel()

	// Output:
	// ......
}
