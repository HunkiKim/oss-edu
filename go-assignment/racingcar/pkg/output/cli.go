package output

import (
	"fmt"
	"racing-car/racingcar/pkg"
	"strings"
)

type CliOutput struct{}

func (co CliOutput) PrintRank(users []*pkg.User) {
	sortUsers(users)

	winners, err := parseWinners(users)
	if err != nil {
		return
	}
	co.printWinners(winners)

	topUsers := users[:min(MaxRank, len(users))]

	for idx, user := range topUsers {
		fmt.Printf("(%d등)%s:%s\n", idx+1, parseNameByLength(user), strings.Repeat("-", user.NumberOfTurns))
	}
}

func (co CliOutput) printWinners(winners []pkg.User) {
	fmt.Print("우승자: ")
	for _, user := range winners {
		fmt.Print(user.Name + " ")
	}
	fmt.Println()
}
