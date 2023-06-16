package json

import (
	"encoding/json"
	"fmt"
	"racing-car/racingcar/pkg/interfaces"
	"racing-car/racingcar/pkg/user"
)

type Writer struct{}

func NewWriter() interfaces.Writer {
	return &Writer{}
}

type output struct {
	Winners  []string     `json:"winners"`
	TopUsers []*user.User `json:"top_users"`
}

func (w *Writer) Write(winners, topUsers []*user.User) error {
	o := output{}
	for _, winner := range winners {
		o.Winners = append(o.Winners, winner.Name)
	}

	o.TopUsers = topUsers

	marshal, err := json.Marshal(o)
	if err != nil {
		return err
	}
	fmt.Println(string(marshal))
	return nil
}
