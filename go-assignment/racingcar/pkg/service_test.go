package pkg

import (
	"testing"
)

func Test_createUsers(t *testing.T) {
	data := []struct {
		testName string
		names    []string
		expected []User
		errMsg   string
	}{{
		"정상 테스트",
		[]string{"user1", "user2", "user3"},
		[]User{
			{"user1", 0},
			{"user2", 0},
			{"user3", 0}},
		""},
		{"중복입력 테스트",
			[]string{"hunki", "hunki", "hunkis", "hunkis", "hunkiss"},
			[]User{
				{name: "hunki", numberOfTurns: 0},
				{name: "hunkis", numberOfTurns: 0},
				{name: "hunkiss", numberOfTurns: 0}},
			""}}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			result := CreateUsers(d.names)
			if len(result) != len(d.expected) {
				t.Error("The number of users created is different")
			}
			// 비즈니스로직
		})
	}
}
