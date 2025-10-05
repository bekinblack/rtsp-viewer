package main

import (
	"stream-viewer/logger"
	"stream-viewer/model"
	"stream-viewer/stream"
	"stream-viewer/ui"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	width  = 640
	height = 480
)

func main() {
	a := app.New()
	w := a.NewWindow("RTSP Stream Viewer")
	w.Resize(fyne.NewSize(width*2, 0))
	w.SetFixedSize(true)

	form := model.Form{}
	form.ChangeIP("1.1.1.1")
	form.ChangePort("1")
	form.ChangeLogin("1")
	form.ChangePassword("1")
	form.ChangeUriHigh("rtsp://test:test@localhost:8554/stream1")
	form.ChangeUriLow("rtsp://test:test@localhost:8554/stream2")

	var streamHigh *stream.Stream
	var streamLow *stream.Stream
	var disconnect *widget.Button
	var check *widget.Button
	var connect *widget.Button

	uriLow := ui.Entry("RTSP URI Low", form.ChangeUriLow)
	uriHigh := ui.Entry("RTSP URI High", form.ChangeUriHigh)

	ip := ui.Entry("IP", form.ChangeIP)
	port := ui.Entry("PORT", form.ChangePort)
	login := ui.Entry("Login", form.ChangeLogin)
	pass := ui.PwdEntry("Password", form.ChangePassword)

	previewHigh := ui.NewViewer(width, height)
	previewLow := ui.NewViewer(width, height)
	statusHigh := ui.Label("Status High")
	statusLow := ui.Label("Status Low")

	messages := widget.NewLabel("Сообщения")

	connect = widget.NewButton("Подключиться", func() {
		check.Disable()

		var wg sync.WaitGroup
		wg.Go(func() {
			var err error
			streamHigh, err = stream.NewStream(form.UriHigh(), width, height)
			if err != nil {
				fyne.Do(func() { statusHigh.SetText("Стрим High:" + err.Error()) })
				return
			}
			fyne.Do(func() { statusHigh.SetText("Стрим High запущен") })

			go logger.LogStream(streamHigh.Err, statusHigh)
			go previewHigh.View(streamHigh.Out)
		})

		wg.Go(func() {
			var err error
			streamLow, err = stream.NewStream(form.UriLow(), width, height)
			if err != nil {
				fyne.Do(func() { statusLow.SetText("Стрим Low:" + err.Error()) })
				return
			}
			fyne.Do(func() { statusLow.SetText("Стрим Low запущен") })

			go logger.LogStream(streamLow.Err, statusLow)
			go previewLow.View(streamLow.Out)
		})

		wg.Wait()
		connect.Disable()
	})

	// подключение недоступно до проверки
	connect.Disable()

	disconnect = widget.NewButton("Отключится", func() {
		err := streamHigh.Close()
		if err != nil {
			statusHigh.SetText(err.Error())
			return
		}
		statusHigh.SetText("Отключено")

		err = streamLow.Close()
		if err != nil {
			statusLow.SetText(err.Error())
			return
		}
		statusLow.SetText("Отключено")

		check.Enable()
		messages.SetText("Нажмите <Проверить> перед подключением")
	})

	check = widget.NewButton("Проверить", func() {
		err := form.Validate()
		if err != nil {
			connect.Disable()
			messages.SetText(err.Error())
			return
		}
		messages.SetText("Проверка пройдена")
		connect.Enable()
	})

	ipRow := container.NewGridWithColumns(4, ip, port, login, pass)

	uriRow := container.NewGridWithColumns(2,
		container.NewVBox(uriHigh, uriLow),
		check,
	)
	previewRow := container.NewVBox(
		container.NewGridWithColumns(2, previewHigh.Content(), previewLow.Content()),
		container.NewGridWithColumns(2, statusHigh, statusLow),
	)

	buttons := container.NewVBox(connect, disconnect)
	msgRow := container.NewGridWithColumns(2, messages, buttons)

	w.SetContent(container.NewVBox(
		ipRow,
		uriRow,
		previewRow,
		msgRow,
	))

	w.ShowAndRun()
}
