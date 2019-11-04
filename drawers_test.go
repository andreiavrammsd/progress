package progress

import (
	"strings"
	"testing"
)

type writerMock struct {
	data []byte
}

func (w *writerMock) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func TestCharacters_DrawDefault(t *testing.T) {
	characters := &Characters{}
	w := &writerMock{}

	count := 3
	for i := 0; i < count; i++ {
		if err := characters.Draw(w); err != nil {
			t.Errorf("writer error: %s", err)
		}
	}

	actual := string(w.data)
	expected := strings.Repeat(".", 3)

	if actual != expected {
		t.Errorf("got: %s, expected: %s", actual, expected)
	}
}

func TestCharacters_DrawWithCharacters(t *testing.T) {
	characters := &Characters{
		Characters: []byte("-"),
	}
	w := &writerMock{}

	count := 3
	for i := 0; i < count; i++ {
		if err := characters.Draw(w); err != nil {
			t.Errorf("writer error: %s", err)
		}
	}

	actual := string(w.data)
	expected := strings.Repeat(string(characters.Characters), 3)

	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}

func TestSpinner_Draw(t *testing.T) {
	spinner := &Spinner{}
	w := &writerMock{}

	if err := spinner.Draw(w); err != nil {
		t.Errorf("writer error: %s", err)
	}
	if err := spinner.Draw(w); err != nil {
		t.Errorf("writer error: %s", err)
	}
	if err := spinner.Draw(w); err != nil {
		t.Errorf("writer error: %s", err)
	}
	if err := spinner.Draw(w); err != nil {
		t.Errorf("writer error: %s", err)
	}

	actual := string(w.data)

	expected := clear + "|" + "\n"
	expected += clear + "/" + "\n"
	expected += clear + "—" + "\n"
	expected += clear + "\\" + "\n"

	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}

func TestBlink_DrawDefault(t *testing.T) {
	blink := &Blink{}
	w := &writerMock{}

	if err := blink.Draw(w); err != nil {
		t.Error("writer error")
	}
	if err := blink.Draw(w); err != nil {
		t.Error("writer error")
	}

	actual := string(w.data)
	expected := clear + string(blink.Characters)

	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}

func TestBlink_DrawWithCharacters(t *testing.T) {
	blink := &Blink{
		Characters: []byte("!!!"),
	}
	w := &writerMock{}

	if err := blink.Draw(w); err != nil {
		t.Error("writer error")
	}
	if err := blink.Draw(w); err != nil {
		t.Error("writer error")
	}

	actual := string(w.data)
	expected := clear + string(blink.Characters)

	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}

func TestArrow_Draw(t *testing.T) {
	arrow := &Arrow{}
	w := &writerMock{}

	if err := arrow.Draw(w); err != nil {
		t.Error("writer error")
	}
	if err := arrow.Draw(w); err != nil {
		t.Error("writer error")
	}

	actual := string(w.data)
	expected := "=>\033[1D=>\033[1D"

	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}
