package progress

import (
	"fmt"
	"testing"
)

type writerMock struct {
	data []byte
}

func (w *writerMock) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func TestBlink_Draw(t *testing.T) {
	blink := &Blink{}
	w := &writerMock{}

	if err := blink.Draw(w); err != nil {
		t.Error("writer error")
	}
	if err := blink.Draw(w); err != nil {
		t.Error("writer error")
	}

	actual := string(w.data)
	expected := clear + "*"

	if actual != expected {
		t.Error(fmt.Sprintf("got: %s, expected: %s", w.data, expected))
	}
}

func TestArrow_Draw(t *testing.T) {
	arrow := &Arrow{}
	w := &writerMock{}

	arrow.Draw(w)
	arrow.Draw(w)
	arrow.Draw(w)

	if err := arrow.Draw(w); err != nil {
		t.Error("writer error")
	}
	if err := arrow.Draw(w); err != nil {
		t.Error("writer error")
	}

	fmt.Println(string(w.data))

	if string(w.data) != clear {
		t.Error(fmt.Sprintf("got: %s, expected: %s", w.data, clear))
	}
}
