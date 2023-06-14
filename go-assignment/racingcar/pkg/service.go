package pkg

import (
	"math/rand"
)

func (user *User) updateNumberOfTurns(turns int) {
	user.NumberOfTurns = rand.Intn(turns + 1)
}

func DoRace(users []*User, turns int) {
	for _, user := range users {
		user.updateNumberOfTurns(turns)
	}
}

func CreateUsers(names []string) []*User {
	users := make([]*User, 0, len(names))
	var duplicated = map[string]bool{}

	for _, name := range names {
		if !duplicated[name] {
			duplicated[name] = true

			users = append(users, &User{
				Name:          name,
				NumberOfTurns: 0,
			})
		}
	}
	return users
}
