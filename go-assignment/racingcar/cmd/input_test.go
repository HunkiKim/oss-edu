package main

import (
	"github.com/google/go-cmp/cmp"
	"racing-car/pkg/user"
	"testing"
)

func Test_inputName(t *testing.T) {
	data := []struct {
		testName string
		name     string
		expected []*user.User
		errMsg   string
	}{
		{
			testName: "정상 테스트",
			name:     "hunki,hunkis,hunkiss",
			expected: []*user.User{{Name: "hunki"}, {Name: "hunkis"}, {Name: "hunkiss"}},
		},
		{
			testName: "경계값 0 테스트",
			name:     ",,",
			errMsg:   "name must be greater than or equal to 1",
		},
		{
			testName: "경계값 11 테스트",
			name:     "hunkihunkih",
			errMsg:   "name must be less than or equal to 10",
		},
		{
			testName: "영어 이외 테스트",
			name:     "hun1",
			errMsg:   "name must be english",
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			users, err := user.CreateUsers(d.name)
			if err != nil && err.Error() != d.errMsg {
				t.Errorf("Expected error message %s, got %s", d.errMsg, err.Error())
			}
			if diff := cmp.Diff(users, d.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
