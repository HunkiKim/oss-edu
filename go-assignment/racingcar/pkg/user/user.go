package user

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
)

var re, _ = regexp.Compile(`[^a-zA-Z]`)

type User struct {
	Name          string
	NumberOfTurns int
}

func DoRace(users []*User, turns int) {
	for _, u := range users {
		updateNumberOfTurns(u, turns)
	}
}

func updateNumberOfTurns(user *User, turns int) {
	user.NumberOfTurns = rand.Intn(turns + 1)
}

func CreateUsers(names string) ([]*User, error) {
	convertedNames, err := convertSlice(names)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(convertedNames))
	var duplicated = map[string]bool{}

	for _, name := range convertedNames {
		if !duplicated[name] {
			duplicated[name] = true

			users = append(users, &User{
				Name:          name,
				NumberOfTurns: 0,
			})
		}
	}
	return users, nil
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
