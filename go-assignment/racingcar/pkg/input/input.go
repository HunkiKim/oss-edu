package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"racing-car/racingcar/pkg"
	"regexp"
	"strings"
)

var (
	re, _  = regexp.Compile(`[^a-zA-Z]`)
	reader = bufio.NewReader(os.Stdin)
)

type Input interface {
	InputNames() ([]string, error)
	InputTurns() (int, error)
}

func CreateInput(ioType pkg.IoType) (Input, error) {
	switch {
	case pkg.Cli == ioType:
		return CliInput{}, nil
	case pkg.File == ioType:
		fmt.Print("파일경로: ")
		texts, err := readFile()
		if err != nil {
			return nil, err
		}
		return FileInput{texts: texts}, nil
	case pkg.Json == ioType:
		fmt.Print("입력:")
		jsonInput, err := inputJsonFormat()
		if err != nil {
			return nil, err
		}
		return jsonInput, nil
	default:
		fmt.Println("err")
		return nil, errors.New("wrong ioType")
	}
}

func convertSlice(input string) ([]string, error) {
	names := strings.Split(input, ",")

	names = deleteZeroString(names)

	if err := verify(names); err != nil {
		return nil, err
	}
	return names, nil
}

func deleteZeroString(names []string) []string {
	newNames := make([]string, 0, len(names))
	for _, name := range names {
		if name != "" {
			newNames = append(newNames, name)
		}
	}
	return newNames
}

func verify(names []string) error {
	for _, name := range names {
		switch {
		case 1 > len(name):
			return errors.New("name must be greater than or equal to 1")
		case 10 < len(name):
			return errors.New("name must be less than or equal to 10")
		case re.MatchString(name):
			return errors.New("name must be english")
		}
	}
	return nil
}
