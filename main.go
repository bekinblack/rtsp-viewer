package main

import (
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

	connectBtn := widget.NewButton("Подключиться", func() {
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

		go logger.LogStdErr("High", streamHigh.Err)
		go logger.LogStdErr("Low", streamLow.Err)

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
