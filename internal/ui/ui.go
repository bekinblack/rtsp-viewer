package ui

import (
	"context"
	"io"
	"rtsp-viewer/internal/config"
	"rtsp-viewer/internal/logger"
	"rtsp-viewer/internal/model"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	fyneApp fyne.App
	ip      *widget.Entry
	port    *widget.Entry

	login *widget.Entry
	pass  *widget.Entry

	pathHigh *widget.Entry
	pathLow  *widget.Entry

	viewHigh *Viewer
	viewLow  *Viewer

	uriHigh    *widget.Label
	uriLow     *widget.Label
	statusHigh *widget.Label
	statusLow  *widget.Label
	messages   *widget.Label

	check      *widget.Button
	connect    *widget.Button
	disconnect *widget.Button
}

const (
	width  = 640
	height = 480
)

type streamer func(ctx context.Context, uri string, width, height int, log *logger.Logger) (Out, Err io.Reader, err error)

func New(ctx context.Context, form model.Form, stream streamer, log *logger.Logger) *App {
	ctx, cancel := context.WithCancel(ctx)
	a := &App{
		fyneApp:  app.New(),
		ip:       Entry(form.IP, "IP", form.SetIP),
		port:     Entry(form.Port, "PORT", form.SetPort),
		login:    Entry(form.Login, "Login", form.SetLogin),
		pass:     PwdEntry(form.Password, "Password", form.SetPassword),
		pathHigh: Entry(form.PathHigh, "URI High Path", form.SetPathHigh),
		pathLow:  Entry(form.PathLow, "URI Low Path", form.SetPathLow),

		uriHigh:  widget.NewLabel("rtsp://{login}:{password}@{ip}:{port}/{path}"),
		uriLow:   widget.NewLabel("rtsp://{login}:{password}@{ip}:{port}/{path}"),
		viewHigh: NewViewer(width, height),
		viewLow:  NewViewer(width, height),

		statusHigh: widget.NewLabel("Остановлен"),
		statusLow:  widget.NewLabel("Остановлен"),
		messages:   widget.NewLabel("Сообщения"),

		check:      widget.NewButton("Проверить", func() {}),
		connect:    widget.NewButton("Подключиться", func() {}),
		disconnect: widget.NewButton("Отключиться", func() {}),
	}

	a.connect.Disable()
	a.disconnect.Disable()

	a.check.OnTapped = func() {
		hi, lo, err := form.Validate()
		if err != nil {
			a.messages.SetText(err.Error())
			a.check.Enable()
			return
		}

		a.disableForm()
		a.uriHigh.SetText(hi)
		a.uriLow.SetText(lo)
		a.messages.SetText("Проверка пройдена")
		a.connect.Enable()

		if err = config.Save(form); err != nil {
			log.Println(err)
		}
	}

	a.disconnect.OnTapped = func() {
		cancel()
		ctx, cancel = context.WithCancel(context.Background())
		a.statusHigh.SetText("Остановлен")
		a.statusLow.SetText("Остановлен")
		a.enableForm()
		a.disconnect.Disable()
	}

	a.connect.OnTapped = func() {
		var wg sync.WaitGroup
		wg.Go(func() {
			streamOut, streamErr, err := stream(ctx, form.UriHigh(), width, height, log)
			if err != nil {
				fyne.Do(func() { a.statusHigh.SetText("Стрим High: " + err.Error()) })
				return
			}
			fyne.Do(func() { a.statusHigh.SetText("Запущен") })

			go log.LogStream(streamErr, a.statusHigh)
			go a.viewHigh.View(streamOut)
		})

		wg.Go(func() {
			streamOut, streamErr, err := stream(ctx, form.UriLow(), width, height, log)
			if err != nil {
				fyne.Do(func() { a.statusLow.SetText("Стрим Low:" + err.Error()) })
				return
			}
			fyne.Do(func() { a.statusLow.SetText("Запущен") })

			go log.LogStream(streamErr, a.statusLow)
			go a.viewLow.View(streamOut)
		})

		wg.Wait()
		a.connect.Disable()
		a.disconnect.Enable()
	}

	a.layout()

	return a
}

func (a *App) Run() {
	a.fyneApp.Run()
}

func (a *App) layout() {
	w := a.fyneApp.NewWindow("RTSP Stream Viewer")
	w.SetFixedSize(true)

	ipRow := container.NewGridWithColumns(4, a.ip, a.port, a.login, a.pass)

	uriRow := container.NewVBox(
		container.NewGridWithColumns(2, a.pathHigh, a.pathLow),
		container.NewGridWithColumns(2, container.NewHScroll(a.uriHigh), container.NewHScroll(a.uriLow)),
	)

	previewRow := container.NewVBox(
		container.NewGridWithColumns(2, a.viewHigh.Content(), a.viewLow.Content()),
		container.NewGridWithColumns(2, a.statusHigh, a.statusLow),
	)

	buttons := container.NewVBox(a.connect, a.disconnect)
	msgRow := container.NewGridWithColumns(2, a.messages, buttons)

	w.SetContent(container.NewVBox(
		ipRow,
		uriRow,
		a.check,
		previewRow,
		msgRow,
	))

	w.Show()
}

func (a *App) enableForm() {
	a.ip.Enable()
	a.port.Enable()
	a.login.Enable()
	a.pass.Enable()
	a.pathHigh.Enable()
	a.pathLow.Enable()
}

func (a *App) disableForm() {
	a.ip.Disable()
	a.port.Disable()
	a.login.Disable()
	a.pass.Disable()
	a.pathHigh.Disable()
	a.pathLow.Disable()
}
