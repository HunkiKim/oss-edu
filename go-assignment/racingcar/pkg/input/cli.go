package input

import (
	"errors"
	"fmt"
)

type CliInput struct{}

func (ci CliInput) InputNames() ([]string, error) {
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

func (ci CliInput) InputTurns() (int, error) {
	var turns int

	fmt.Print("주행 횟수:")
	_, err := fmt.Scan(&turns)
	if err != nil {
		return 0, err
	}
	if turns < 1 {
		return 0, errors.New("turns cannot be less than one")
	}
	return turns, err
}
