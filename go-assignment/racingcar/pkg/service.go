package pkg

import (
	"math/rand"
	"sort"
)

type User struct {
	name          string
	numberOfTurns int
}

func CountNumberOfTurns(names []string, numberOfTurns int) []User { // 이름 수정
	users := createUsers(names)

	for idx := range users {
		users[idx].countNumberOfTurns(numberOfTurns)
	}

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
	return users
}

func createUsers(names []string) []User {
	users := make([]User, 0, len(names))
	var duplicatedMap = map[string]bool{}

	for _, name := range names { // 일급함수로 바꿔보기
		if duplicatedMap[name] {
			continue
		}
		duplicatedMap[name] = true

		users = append(users, User{
			name:          name,
			numberOfTurns: 0,
		})
	}

	return users
}

func (user *User) countNumberOfTurns(numberOfTurns int) {
	user.numberOfTurns = rand.Intn(numberOfTurns + 1)
}
