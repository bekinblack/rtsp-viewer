package main

import (
	"rtsp-viewer/internal/config"
	"rtsp-viewer/internal/logger"
	"rtsp-viewer/internal/stream"
	"rtsp-viewer/internal/ui"

	"golang.org/x/net/context"
)

func main() {
	log := logger.New()
	defer log.Close()

	form, err := config.Load()
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()
	app := ui.New(ctx, form, stream.New, log)
	app.Run()
}
