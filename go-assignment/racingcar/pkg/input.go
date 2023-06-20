package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var re, _ = regexp.Compile(`[^a-zA-Z]`)

func InputNames() ([]string, error) {
	var names string

	fmt.Print("이름:")
	_, err := fmt.Scan(&names)
	if err != nil {
		return nil, err
	}

	convertedNames, err := convertSlice(names)
	if err != nil {
		return nil, err
	}
	return convertedNames, nil
}

func InputTurns() (int, error) {
	var turns int

	fmt.Print("도는 횟수:")
	_, err := fmt.Scan(&turns)
	if err != nil {
		return 0, err
	}
	if turns < 1 {
		return 0, errors.New("turns cannot be less than one")
	}
	return turns, err
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
