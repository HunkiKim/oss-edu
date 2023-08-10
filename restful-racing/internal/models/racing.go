package models

type Racing struct {
	MaxTurns int `json:"max_turns"`
}

var (
	Racings  map[int64]Racing = map[int64]Racing{}
	racingId int64            = 0
)

func AddRacing(r Racing) {
	Racings[racingId] = r
	racingId += 1
}
