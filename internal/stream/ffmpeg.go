package stream

import (
	"context"
	"fmt"
	"os/exec"
)

func ffmpegCmd(ctx context.Context, uri string, width, height int) *exec.Cmd {
	return exec.CommandContext(ctx, "ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", uri,
		"-loglevel", "error",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
		"-vsync", "0",
		"-f", "rawvideo",
		"-pix_fmt", "rgba",
		"pipe:1",
	)
}
