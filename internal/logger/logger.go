package logger

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
)

type textSetter interface {
	SetText(s string)
}

type Logger struct {
	file *os.File
	log  *log.Logger
}

func New() *Logger {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return &Logger{
		file: file,
		log:  log.New(file, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l Logger) Close() {
	l.file.Close()
}

func (l Logger) Println(v ...any) {
	l.log.Println(v...)
}

func (l Logger) LogStream(from io.Reader, to textSetter) {
	go func() {
		scanner := bufio.NewScanner(from)
		for scanner.Scan() {
			text := scanner.Text()
			switch {
			case strings.Contains(text, "401 Unauthorized"):
				fyne.Do(func() { to.SetText("Ошибка: проверьте логин и пароль") })
			case strings.Contains(text, "403 Forbidden"):
				fyne.Do(func() { to.SetText("Ошибка: доступ запрещен") })
			case strings.Contains(text, "Error opening input"):
				fyne.Do(func() { to.SetText("Ошибка: проверьте адрес стрима") })
			case strings.Contains(text, "Unrecognized option"):
				fyne.Do(func() { to.SetText("Ошибка: проверьте опции ffmpeg") })
			default:
				l.log.Println(text)
			}
		}
	}()
}
