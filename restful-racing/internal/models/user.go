package models

type User struct {
	Name     string `json:"name"`
	Turns    int    `json:"turns"`
	RacingId int64  `json:"racing_id"`
}

var users map[int64]User = map[int64]User{}
