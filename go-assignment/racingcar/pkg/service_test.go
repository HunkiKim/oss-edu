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
	}{
		{
			testName: "정상 테스트",
			names:    []string{"user1", "user2", "user3"},
			expected: []User{
				{"user1", 0},
				{"user2", 0},
				{"user3", 0},
			},
		},
		{
			testName: "중복입력 테스트",
			names:    []string{"hunki", "hunki", "hunkis", "hunkis", "hunkiss"},
			expected: []User{
				{Name: "hunki", NumberOfTurns: 0},
				{Name: "hunkis", NumberOfTurns: 0},
				{Name: "hunkiss", NumberOfTurns: 0},
			},
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			result := CreateUsers(d.names)
			if len(result) != len(d.expected) {
				t.Error("The number of users created is different")
			}
		})
	}
}
