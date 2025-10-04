package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Logger struct {
	label *widget.Label
}

func NewLogger(label string) *Logger {
	return &Logger{label: widget.NewLabel(label)}
}

func (l *Logger) Content() fyne.CanvasObject {
	return l.label
}

func (l *Logger) Log(msgs ...any) {
	fyne.Do(func() {
		l.label.SetText(fmt.Sprintf("%v", msgs...))
	})
}
