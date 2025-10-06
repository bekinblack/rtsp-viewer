package stream

import (
	"errors"
	"io"
	"os/exec"
)

type Stream struct {
	cmd *exec.Cmd
	Out io.Reader
	Err io.Reader
}

func NewStream(uri string, width, height int) (*Stream, error) {
	cmd := NewFfmpeg(uri, width, height)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	s := &Stream{
		cmd: cmd,
		Out: stdout,
		Err: stderr,
	}

	if err := s.cmd.Start(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Stream) Close() error {
	if s != nil {
		return s.cmd.Process.Kill()
	}
	return errors.New("Стрим не запущен")
}
