package output

import (
	"fmt"
	"os"
	"racing-car/racingcar/pkg"
	"strings"
)

type FileOutput struct{}

const filePath = "./result.txt"

func (fo FileOutput) PrintRank(users []*pkg.User) {
	sortUsers(users)

	winners, err := parseWinners(users)
	if err != nil {
		return
	}
	winnersPrint := fo.sprintWinners(winners)

	createdFile, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer func(createdFile *os.File) {
		err := createdFile.Close()
		if err != nil {

		}
	}(createdFile)

	_, err = fmt.Fprintf(createdFile, winnersPrint)
	if err != nil {
		return
	}

	topUsers := users[:min(MaxRank, len(users))]
	for idx, user := range topUsers {
		_, err := fmt.Fprintf(createdFile, fmt.Sprintf("(%d등)%s:%s\n", idx+1, parseNameByLength(user), strings.Repeat("-", user.NumberOfTurns)))
		if err != nil {
			return
		}
	}
}

func (fo FileOutput) sprintWinners(winners []pkg.User) string {
	result := fmt.Sprint("우승자: ")
	for _, user := range winners {
		result += fmt.Sprint(user.Name + " ")
	}
	result += fmt.Sprintln()
	return result
}
