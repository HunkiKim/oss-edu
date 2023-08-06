package interfaces

import (
	"racing-car/pkg/user"
)

type Reader interface {
	Read() (string, int, error)
}

type Writer interface {
	Write(winners, topUsers []*user.User) error
}
