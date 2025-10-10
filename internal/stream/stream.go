package stream

import (
	"context"
	"io"
	"rtsp-viewer/internal/logger"
)

func New(ctx context.Context, uri string, width, height int, log *logger.Logger) (Out, Err io.Reader, err error) {
	cmd := ffmpegCmd(ctx, uri, width, height)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Println(err)
			return
		}
	}()

	return stdout, stderr, nil
}
