package components

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

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

// LogStdErr сканирует stderr и логирует только нужные сообщения.
func (l *Logger) LogStdErr(streamLabel string, r io.Reader) {
	scanner := bufio.NewScanner(r)

	// Регулярные выражения для поиска
	authRe := regexp.MustCompile(`(?i)(401|403)`)
	timeoutRe := regexp.MustCompile(`(?i)(timeout|timed out)`)
	hostRe := regexp.MustCompile(`(?i)(connection refused|no route to host|host unreachable)`)
	codecRe := regexp.MustCompile(`(?i)(unknown codec|could not find codec|decoder.*not found)`)

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case authRe.MatchString(line):
			l.Log(streamLabel + ": Ошибка авторизации (401/403) при подключении")
		case timeoutRe.MatchString(line):
			l.Log(streamLabel + ": Таймаут подключения к RTSP-потоку")
		case hostRe.MatchString(line):
			l.Log(streamLabel + ": Хост недоступен или сетевой сбой")
		case codecRe.MatchString(line):
			l.Log(streamLabel + ": Не удалось декодировать поток — проверьте установку необходимых кодеков/библиотек")
		}
	}
	if err := scanner.Err(); err != nil {
		l.Log(streamLabel + ": Ошибка чтения stderr: " + err.Error())
	}
}
