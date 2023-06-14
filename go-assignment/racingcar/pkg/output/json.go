package output

import (
	"encoding/json"
	"fmt"
	"racing-car/racingcar/pkg"
)

const (
	winnersKey = "winners"
)

type JsonOutput struct {
	Winners  []string    `json:"winners"`
	TopUsers []*pkg.User `json:"top_users"`
}

func (jo JsonOutput) PrintRank(users []*pkg.User) error {
	sortUsers(users)

	winners, err := parseWinners(users)
	if err != nil {
		return err
	}
	for _, winner := range winners {
		jo.Winners = append(jo.Winners, winner.Name)
	}

	jo.TopUsers = users[:min(MaxRank, len(users))]

	marshal, err := json.Marshal(jo)
	if err != nil {
		return err
	}
	fmt.Println(string(marshal))

	return nil
}
