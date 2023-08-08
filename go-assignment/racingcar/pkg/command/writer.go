package command

import (
	"fmt"
	"racing-car/pkg/interfaces"
	"racing-car/pkg/user"
	"strings"
)

type Writer struct{}

func NewWriter() interfaces.Writer {
	return &Writer{}
}

func (w *Writer) Write(winners, topUsers []*user.User) error {
	w.printWinners(winners)

	for idx, u := range topUsers {
		fmt.Printf("(%d등)%s:%s\n", idx+1, u.Name, strings.Repeat("-", u.NumberOfTurns))
	}

	return nil
}

func (w *Writer) printWinners(winners []*user.User) {
	fmt.Print("우승자: ")
	for _, u := range winners {
		fmt.Print(u.Name + " ")
	}
	fmt.Println()
}
