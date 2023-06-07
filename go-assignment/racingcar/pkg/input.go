package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func InputAll() ([]string, int, error) { // Init or Input (All)
	var name string
	fmt.Print("이름:")
	fmt.Scan(&name)
	names, err := inputName(name)
	if err != nil {
		return nil, 0, err
	}

	fmt.Print("도는 횟수:")
	var numberOfTurns int // count, turns, 등등..
	_, err = fmt.Scan(&numberOfTurns)
	if err != nil {
		return nil, 0, err
	}
	if numberOfTurns <= 0 {
		return nil, 0, errors.New("numberOfTurns cannot be less than or equal to zero")
	}

	return names, numberOfTurns, err
}

func inputName(name string) ([]string, error) { //시그니쳐를 생각해서 바꿔보기
	names := strings.Split(name, ",")

	for i := 0; i < len(names); i++ {
		if checkErr := checkName(string(names[i])); checkErr != "" { // 에러를 바로 반환하도록 수정
			return nil, errors.New(checkErr)
		}
	}
	return names, nil
}

func checkName(name string) string {
	switch {
	case 1 > len(name):
		return "name must be greater than or equal to 1"
	case 10 < len(name):
		return "name must be less than or equal to 10"
	}
	if matched, _ := regexp.MatchString(`[^ㄱ-ㅎㅏ-ㅣ가-힣a-zA-Z]`, name); matched {
		return "name must be korean or english"
	}
	return ""
}
