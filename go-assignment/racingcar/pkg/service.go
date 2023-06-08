package pkg

import (
	"math/rand"
)

type user struct {
	name          string
	numberOfTurns int
}

func (user *user) countNumberOfTurns(turns int) {
	user.numberOfTurns = rand.Intn(turns + 1)
}

func DoRace(users []*user, turns int) {
	for _, user := range users {
		user.countNumberOfTurns(turns)
	}
}

func CreateUsers(names []string) []*user {
	users := make([]*user, 0, len(names))
	var duplicated = map[string]bool{}

	for _, name := range names {
		if !duplicated[name] {
			duplicated[name] = true

			users = append(users, &user{
				name:          name,
				numberOfTurns: 0,
			})
		}
	}
	return users
}
