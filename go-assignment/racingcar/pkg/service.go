package pkg

import (
	"math/rand"
	"sort"
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

func SortUsers(users []*User) {
	sort.Slice(users, func(i, j int) bool {
		switch {
		case users[i].numberOfTurns > users[j].numberOfTurns:
			return true
		case users[i].numberOfTurns < users[j].numberOfTurns:
			return false
		case users[i].name < users[j].name:
			return true
		default:
			return false
		}
	})
}

func (user *User) countNumberOfTurns(turns int) {
	user.numberOfTurns = rand.Intn(turns + 1)
}
