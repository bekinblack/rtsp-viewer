package stream

import (
	"io"
	"os/exec"
)

type Stream struct {
	cmd *exec.Cmd
	Out io.Reader
	Err io.Reader
}

func NewStream(cmd *exec.Cmd) (*Stream, error) {
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
	return s.cmd.Process.Kill()
}
