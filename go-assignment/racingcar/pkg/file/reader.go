package file

import (
	"errors"
	"fmt"
	"os"
	"racing-car/pkg/interfaces"
	"strconv"
	"strings"
)

type Reader struct{}

func NewReader() interfaces.Reader {
	return &Reader{}
}

func (r *Reader) Read() (string, int, error) {
	fmt.Print("파일경로: ")
	texts, err := r.readFile()
	if err != nil {
		return "", -1, err
	}

	names, err := r.parseName(texts[0])
	if err != nil {
		return "", 0, err
	}

	turns, err := r.parseTurns(texts[1])
	if err != nil {
		return "", 0, err
	}
	return names, turns, nil
}

func (r *Reader) readFile() ([]string, error) {
	var filePath string
	_, err := fmt.Scan(&filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	texts := strings.Split(string(file), "\n")
	if len(texts) != 2 {
		return nil, errors.New("file is not two lines")
	}
	return texts, nil
}

func (r *Reader) parseTurns(text string) (int, error) {
	parsedText := strings.Split(text, ":")
	if len(parsedText) != 2 || parsedText[0] != "주행 횟수" {
		return 0, errors.New("file invalid with turns line ")
	}

	turns, err := strconv.Atoi(parsedText[1])
	if err != nil {
		return 0, err
	}
	return turns, nil
}

func (r *Reader) parseName(text string) (string, error) {
	parsedText := strings.Split(text, ":")
	if len(parsedText) != 2 || parsedText[0] != "이름" {
		return "", errors.New("file invalid with name line ")
	}
	return parsedText[1], nil
}
