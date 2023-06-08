package pkg

import (
	"math/rand"
)

type User struct {
	name          string
	numberOfTurns int
}

func DoRace(users []*User, turns int) {
	for _, user := range users {
		user.countNumberOfTurns(turns)
	}
}

func CreateUsers(names []string) []*User {
	users := make([]*User, 0, len(names))
	var duplicatedMap = map[string]bool{}

	for _, name := range names {
		if !duplicatedMap[name] {
			duplicatedMap[name] = true

			users = append(users, &User{
				name:          name,
				numberOfTurns: 0,
			})
		}
	}

	return users
}

func (user *User) countNumberOfTurns(turns int) {
	user.numberOfTurns = rand.Intn(turns + 1)
}
