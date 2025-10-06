package stream

import (
	"fmt"
	"os/exec"
)

func NewFfmpeg(uri string, width, height int) *exec.Cmd {
	return exec.Command("ffmpeg",
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
