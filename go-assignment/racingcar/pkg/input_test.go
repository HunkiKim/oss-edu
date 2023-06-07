package pkg

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_inputName(t *testing.T) {
	data := []struct {
		testName string
		name     string
		expected []string
		errMsg   string
	}{
		{"정상 테스트", "hunki,hunkis,hunkiss", []string{"hunki", "hunkis", "hunkiss"}, ""},
		{"경계값 0 테스트", "hunki,,", nil, "name must be greater than or equal to 1"},
		{"경계값 11 테스트", "hunkihunkih", nil, "name must be less than or equal to 10"},
		{"한글영어 이외 테스트", "hun1", nil, "name must be korean or english"},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			result, err := getNamesByInput(d.name)
			if err != nil && err.Error() != d.errMsg {
				t.Errorf("Expected error message %s, got %s", d.errMsg, err.Error())
			}
			if diff := cmp.Diff(result, d.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
