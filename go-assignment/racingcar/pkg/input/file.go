package input

import (
	"errors"
	"strconv"
	"strings"
)

type FileInput struct {
	texts []string
}

const (
	nameFormat  = "이름"
	turnsFormat = "도는 횟수"
	nameIdx     = 0
	turnsIdx    = 1
)

func (fi FileInput) InputNames() ([]string, error) {
	names, err := parseName(fi.texts[nameIdx])
	if err != nil {
		return nil, err
	}

	convertedNames, err := convertSlice(names)
	if err != nil {
		return nil, err
	}
	return convertedNames, nil
}

func (fi FileInput) InputTurns() (int, error) {
	parsedText := strings.Split(fi.texts[turnsIdx], ":")
	if len(parsedText) != 2 || parsedText[0] != turnsFormat {
		return 0, errors.New("file invalid with turns line ")
	}

	turns, err := strconv.Atoi(parsedText[1])
	if err != nil {
		return 0, err
	}
	return turns, nil
}

func parseName(text string) (string, error) {
	parsedText := strings.Split(text, ":")
	if len(parsedText) != 2 || parsedText[0] != nameFormat {
		return "", errors.New("file invalid with name line ")
	}
	return parsedText[1], nil
}
