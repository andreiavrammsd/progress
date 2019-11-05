package progress

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	p := New(nil)

	// Default Drawer
	_, ok := p.drawer.(*Spinner)
	if !ok {
		t.Error("expected default Drawer: Spinner")
	}

	// Default Writer
	if p.writer != os.Stdout {
		t.Error("expected default Drawer: os.Stdout")
	}

	// Error callback
	if p.error != nil {
		t.Error("expected no error callback")
	}

	if p.ctx == nil {
		t.Error("expected context")
	}
}

type drawerErrorMock struct {
}

func (d *drawerErrorMock) Draw(io.Writer) (err error) {
	err = errors.New("writer error at drawer")
	return
}

func TestProgress_ProgressWithError(t *testing.T) {
	c := &Config{}

	c.Drawer = &drawerErrorMock{}

	var drawerError error
	c.Error = func(err error) {
		drawerError = err
	}

	var cancel func()
	c.Context, cancel = context.WithCancel(context.Background())

	p := New(c)
	p.Progress()
	cancel()

	if drawerError == nil {
		t.Error("error expected")
	}
}

type drawerMock struct {
	char string
}

func (d *drawerMock) Draw(w io.Writer) (err error) {
	_, err = w.Write([]byte(d.char))
	return
}

func TestProgress_Progress(t *testing.T) {
	d := &drawerMock{}
	w := &writerMock{}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Config{
		Drawer:  d,
		Writer:  w,
		Context: ctx,
	}
	p := New(c)
	p.Progress()
	cancel()

	actual := string(w.data)
	expected := d.char

	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}
