package input

import (
	"bufio"
	"encoding/json"
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
		switch {
		case len(texts) > 2:
			return nil, errors.New("file line exceeded two lines")
		case len(texts) < 2:
			return nil, errors.New("file line is less than two lines")
		}
		return FileInput{texts: texts}, nil
	case pkg.Json == ioType:
		fmt.Print("입력:")
		var input []byte
		input, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("입력 에러:", err)
			return nil, err
		}

		var jsonInput JsonInput
		err = json.Unmarshal([]byte(input), &jsonInput)
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

	if err := verify(names); err != nil {
		return nil, err
	}
	return names, nil
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
