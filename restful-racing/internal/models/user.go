package models

import "math/rand"

type User struct {
	Name     string `json:"name"`
	Turns    int    `json:"turns"`
	RacingId int64  `json:"racing_id"`
}

var (
	Users        = map[int64]*User{}
	userId int64 = 0
)

func NewUser(name string, maxTurns int, racingId int64) *User {
	return &User{Name: name, Turns: setTurns(maxTurns), RacingId: racingId}
}

func AddUser(user *User) int64 {
	Users[userId] = user
	userId += 1
	return userId - 1
}

func UpdateTurns(id int64, maxTurns int) {
	Users[id].Turns = setTurns(maxTurns)
}

func setTurns(maxTurns int) int {
	return rand.Intn(maxTurns + 1)
}
