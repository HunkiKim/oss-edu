package command

import (
	"fmt"
	"racing-car/racingcar/pkg/user"
	"strings"
)

type Writer struct{}

func (w *Writer) Write(winners, topUsers []*user.User) error {
	w.printWinners(winners)

	for idx, user := range topUsers {
		fmt.Printf("(%d등)%s:%s\n", idx+1, user.Name, strings.Repeat("-", user.NumberOfTurns))
	}

	return nil
}

func (w *Writer) printWinners(winners []*user.User) {
	fmt.Print("우승자: ")
	for _, user := range winners {
		fmt.Print(user.Name + " ")
	}
	fmt.Println()
}
