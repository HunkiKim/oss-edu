package output

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"racing-car/racingcar/pkg"
	"strings"
)

type FileOutput struct{}

const filePath = "./result.txt"

func (fo FileOutput) PrintRank(users []*pkg.User) error {
	sortUsers(users)

	winners, err := parseWinners(users)
	if err != nil {
		return err
	}
	winnersPrint := fo.sprintWinners(winners)

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	createdFile, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer func(createdFile *os.File) {
		err = createdFile.Close()
		if err != nil {
			log.Fatal("파일 닫기 에러")
		}
		fmt.Println("파일이 저장되었습니다. (" + absPath + ")")
	}(createdFile)

	_, err = fmt.Fprintf(createdFile, winnersPrint)
	if err != nil {
		return err
	}

	topUsers := users[:min(MaxRank, len(users))]
	for idx, user := range topUsers {
		_, err := fmt.Fprintf(createdFile, fmt.Sprintf("(%d등)%s:%s\n", idx+1, parseNameByLength(user), strings.Repeat("-", user.NumberOfTurns)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (fo FileOutput) sprintWinners(winners []pkg.User) string {
	result := fmt.Sprint("우승자: ")
	for _, user := range winners {
		result += fmt.Sprint(user.Name + " ")
	}
	result += fmt.Sprintln()
	return result
}
