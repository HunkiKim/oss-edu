package util

import (
	"bufio"
	"errors"
	"os"
	"racing-car/racingcar/pkg/command"
	"racing-car/racingcar/pkg/file"
	"racing-car/racingcar/pkg/json"
	"racing-car/racingcar/pkg/user"
)

type Reader interface {
	Read() (string, int, error)
}

func NewReader(format string) (Reader, error) {
	switch format {
	case "Command":
		return &command.Reader{}, nil
	case "File":
		return &file.Reader{}, nil
	case "Json":
		return &json.Reader{Reader: bufio.NewReader(os.Stdin)}, nil
	default:
		return nil, errors.New("wrong format")
	}
}

type Writer interface {
	Write(winners, topUsers []*user.User) error
}

func NewWriter(format string) (Writer, error) {
	switch format {
	case "Command":
		return &command.Writer{}, nil
	case "File":
		return &file.Writer{}, nil
	case "Json":
		return &json.Writer{}, nil
	default:
		return nil, errors.New("wrong format")
	}
}
