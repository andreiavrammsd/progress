package progress

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func ExampleNew() {
	c := &Config{
		Drawer: &Characters{Characters: []byte(".")},
	}
	p := New(c)

	p.Progress()
	p.Progress()
	p.Progress()

	// Output:
	// ...
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

func ExampleNew2() {
	c := Config{
		Drawer: &Characters{Characters: []byte("+")},
		Writer: os.Stdout,
		Error: func(err error) {
			fmt.Println(err)
		},
	}
	p := New(&c)

	wg := sync.WaitGroup{}
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func(p *Progress) {
			defer wg.Done()
			p.Progress()
		}(p)
	}
	wg.Wait()
	time.Sleep(time.Second)

	// Output:
	// ++++
}
