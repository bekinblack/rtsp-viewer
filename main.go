package main

import (
	"bufio"
	"io"
	"regexp"
	"stream-viewer/components"
	"stream-viewer/model"
	"stream-viewer/stream"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Подсказки
	logger := components.NewLogger("Сообщения:")
	m := model.New(nil, nil, logger)

	a := app.New()
	w := a.NewWindow("RTSP Stream Viewer")
	w.Resize(fyne.NewSize(1280, 0))

	// Верхние поля
	ipEntry := components.NewEntry("IP", m.ChangeIP)
	portEntry := components.NewEntry("Port", m.ChangePort)
	loginEntry := components.NewEntry("Login", m.ChangeLogin)
	passEntry := components.NewPwdEntry("Password", m.ChangePassword)

	// RTSP URI
	rtspLowEntry := components.NewEntry("RTSP URI Low", m.ChangeUriLow)
	rtspHighEntry := components.NewEntry("RTSP URI High", m.ChangeUriHigh)

	checkBtn := widget.NewButton("Проверить", m.Validate)

	// Video area
	viewerHigh := components.NewViewer(640, 480)
	viewerLow := components.NewViewer(320, 240)

	var streamHigh *stream.Stream
	var streamLow *stream.Stream

	connectBtn := widget.NewButton("Подключится", func() {
		urlHigh := "rtsp://test:test@localhost:8554/stream"
		urlLow := "rtsp://test:test@localhost:8554/stream"

		cmdHigh := stream.NewFfmpeg(urlHigh, 640, 480)
		cmdLow := stream.NewFfmpeg(urlLow, 320, 240)

		var err error
		streamHigh, err = stream.NewStream(cmdHigh)
		if err != nil {
			logger.Log(err)
			return
		}

		streamLow, err = stream.NewStream(cmdLow)
		if err != nil {
			logger.Log(err)
			return
		}

		go filterAndLogErrors(streamHigh.Err, logger)
		go filterAndLogErrors(streamLow.Err, logger)

		go viewerHigh.View(streamHigh.Out)
		go viewerLow.View(streamLow.Out)

		logger.Log("Подключено: читаем поток...")

	})

	disconnectBtn := widget.NewButton("Отключится", func() {
		err := streamHigh.Close()
		if err != nil {
			logger.Log("disconnect:", err)
			return
		}
		err = streamLow.Close()
		if err != nil {
			logger.Log("disconnect:", err)
			return
		}

		logger.Log("Отключено")
	})

	// Layout: верхние поля и RTSP
	topRow := container.NewGridWithColumns(4,
		ipEntry, portEntry,
		loginEntry, passEntry,
	)

	rtspRow := container.NewGridWithColumns(2,
		container.NewGridWithRows(2,
			rtspHighEntry,
			rtspLowEntry,
		),
		checkBtn,
	)

	previewRow := container.NewGridWithColumns(2,
		viewerHigh.Content(),
		viewerLow.Content(),
	)

	btnRow := container.NewGridWithColumns(2,
		logger.Content(),
		container.NewGridWithRows(2,
			connectBtn,
			disconnectBtn,
		),
	)

	// Основной layout
	mainContainer := container.NewVBox(
		topRow,
		rtspRow,
		previewRow,
		btnRow,
	)

	w.SetContent(mainContainer)

	w.ShowAndRun()

}

// filterAndLogErrors сканирует stderr и логирует только нужные сообщения.
func filterAndLogErrors(r io.Reader, logger model.Logger) {
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
			logger.Log("Ошибка авторизации (401/403) при подключении")
		case timeoutRe.MatchString(line):
			logger.Log("Таймаут подключения к RTSP-потоку")
		case hostRe.MatchString(line):
			logger.Log("Хост недоступен или сетевой сбой")
		case codecRe.MatchString(line):
			logger.Log("Не удалось декодировать поток — проверьте установку необходимых кодеков/библиотек")
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Log("Ошибка чтения stderr: " + err.Error())
	}
}
