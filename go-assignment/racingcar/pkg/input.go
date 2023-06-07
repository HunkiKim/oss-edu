package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func InputAll() ([]string, int, error) {
	var input string
	fmt.Print("이름:")
	fmt.Scan(&input)
	names, err := getNamesByInput(input)
	if err != nil {
		return nil, 0, err
	}

	fmt.Print("도는 횟수:")
	var turns int
	_, err = fmt.Scan(&turns)
	if err != nil {
		return nil, 0, err
	}
	if turns < 1 {
		return nil, 0, errors.New("turns cannot be less than one")
	}

	return names, turns, err
}

func getNamesByInput(input string) ([]string, error) {
	names := strings.Split(input, ",")

	for _, name := range names {
		if err := checkName(name); err != nil {
			return nil, err
		}
	}
	return names, nil
}

func checkName(name string) error {
	switch matched, _ := regexp.MatchString(`[^a-zA-Z]`, name); {
	case 1 > len(name):
		return errors.New("name must be greater than or equal to 1")
	case 10 < len(name):
		return errors.New("name must be less than or equal to 10")
	case matched:
		return errors.New("name must be korean or english")
	default:
		return nil
	}
}
