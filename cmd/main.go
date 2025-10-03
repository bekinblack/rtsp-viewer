package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/fyne-streamer/video"
)

func main() {
	a := app.New()
	w := a.NewWindow("RTSP-stream")
	w.Resize(fyne.NewSize(900, 600)) // Примерный размер
	w.SetFixedSize(true)

	// Верхние поля
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("IP")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port")

	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Login")

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Password")

	// RTSP URI
	rtspHigh := widget.NewEntry()
	rtspHigh.SetPlaceHolder("RTSP URI 1 (High)")
	rtspLow := widget.NewEntry()
	rtspLow.SetPlaceHolder("RTSP URI 2 (Low)")

	checkBtn := widget.NewButton("Проверить", func() {})

	// Подсказки
	statusMsg := widget.NewLabel("Подсказки/сообщения валидации:")

	// Preview панели
	//highPreview := widget.NewLabel("Preview High\nСтатус: Отключено")
	//lowPreview := widget.NewLabel("Preview Low\nСтатус: Отключено")

	url := "rtsp://localhost:8554/stream"

	lowPreview := video.NewViewer()
	lowPreview.SetMinSize(fyne.NewSize(500, 400))

	highPreview := video.NewViewer()
	highPreview.SetMinSize(fyne.NewSize(500, 400))

	// Кнопки
	connectBtn := widget.NewButton("Подключиться", func() {
		if lowPreview.IsPlaying() || highPreview.IsPlaying() {
			return
		}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			lowPreview.Clear()
			lowPipeline := buildPipeline(url, 10, 10)

			fyne.Do(func() {
				if err := lowPreview.SetPipelineFromString(lowPipeline); err != nil {
					log.Println("SetPipelineFromString failed: ", err)
					return
				}

				err := lowPreview.Play()
				if err != nil {
					log.Println("connection failed: ", err)
					return
				}
			})
		}()

		go func() {
			defer wg.Done()
			highPreview.Clear()
			highPipeline := buildPipeline(url, 10, 80)

			fyne.Do(func() {
				if err := highPreview.SetPipelineFromString(highPipeline); err != nil {
					log.Println("SetPipelineFromString failed: ", err)
					return
				}

				err := highPreview.Play()
				if err != nil {
					log.Println("connection failed: ", err)
					return
				}
			})
		}()

		go func() {
			wg.Wait()
			fyne.Do(func() {
				statusMsg.SetText("Подключено")
			})

		}()
	})

	disconnectBtn := widget.NewButton("Отключиться", func() {
		if !lowPreview.IsPlaying() && !highPreview.IsPlaying() {
			return
		}

		var wg sync.WaitGroup

		if lowPreview.IsPlaying() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := lowPreview.Stop(); err != nil {
					log.Println("low Stop error: ", err)
				}

				if err := lowPreview.Clear(); err != nil {
					log.Printf("low Clear error: %v", err)
				}
			}()
		}

		if highPreview.IsPlaying() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := highPreview.Stop(); err != nil {
					log.Println("high Stop error: ", err)
				}

				if err := highPreview.Clear(); err != nil {
					log.Printf("high Clear error: %v", err)
				}
			}()
		}

		go func() {
			wg.Wait()
			time.Sleep(50 * time.Millisecond)
			fyne.Do(func() {
				statusMsg.SetText("Отключено")
			})
		}()
	})

	// Layout: верхние поля и RTSP
	topRow := container.NewGridWithColumns(4, ipEntry, portEntry, loginEntry, passEntry)
	rtspRow := container.NewGridWithColumns(2, container.NewGridWithRows(2, rtspHigh, rtspLow), checkBtn)
	previewRow := container.NewGridWithColumns(2, highPreview, lowPreview)
	btnRow := container.NewGridWithColumns(2, statusMsg, container.NewGridWithRows(2, connectBtn, disconnectBtn))

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

func buildPipeline(url string, fps, quality int) string {
	return fmt.Sprintf(`
        rtspsrc location=%s latency=200 protocols=tcp !
        rtph264depay !
        decodebin name=dec dec. !
		videorate ! 
		video/x-raw,framerate=%d/1 !
		jpegenc name={{ .ImageEncoderElementName }} quality=%d !
		appsink name={{ .AppSinkElementName }} drop=false sync=true
    `, url, fps, quality)
}
