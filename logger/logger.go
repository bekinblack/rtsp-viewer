package logger

import (
	"bufio"
	"io"
	"log"
	"strings"

	"fyne.io/fyne/v2"
)

type TextSetter interface {
	SetText(s string)
}

func LogStream(from io.Reader, to TextSetter) {
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
				log.Println(text)
			}
		}
	}()
}
