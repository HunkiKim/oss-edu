package command

import (
	"errors"
	"fmt"
)

type Reader struct{}

func (r *Reader) Read() (string, int, error) {
	names, err := r.inputNames()
	if err != nil {
		return "", 0, err
	}

	turns, err := r.inputTurns()
	if err != nil {
		return "", 0, err
	}

	return names, turns, nil
}

func (r *Reader) inputNames() (string, error) {
	var names string

	fmt.Print("이름:")
	_, err := fmt.Scan(&names)
	if err != nil {
		return "", err
	}

	return names, nil
}

func (r *Reader) inputTurns() (int, error) {
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
