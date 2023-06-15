package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"racing-car/racingcar/pkg/user"
)

const filePath = "./result.txt"

type Writer struct{}

func (w *Writer) Write(winners, topUsers []*user.User) error {
	winnersPrint := w.sprintWinners(winners)

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
	}(createdFile)

	_, err = fmt.Fprintf(createdFile, winnersPrint)
	if err != nil {
		return err
	}

	for idx, u := range topUsers {
		_, err := fmt.Fprintf(createdFile, fmt.Sprintf("(%d등)%s:%s\n", idx+1, u.Name, strings.Repeat("-", u.NumberOfTurns)))
		if err != nil {
			return err
		}
	}

	fmt.Println("파일이 저장되었습니다. (" + absPath + ")")

	return nil
}

func (w *Writer) sprintWinners(winners []*user.User) string {
	result := fmt.Sprint("우승자: ")
	for _, u := range winners {
		result += fmt.Sprint(u.Name + " ")
	}
	result += fmt.Sprintln()
	return result
}
