package input

import (
	"encoding/json"
	"fmt"
)

type JsonInput struct {
	Names string
	Turns int
}

func inputJsonFormat() (JsonInput, error) {
	var input []byte
	input, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("입력 에러:", err)
		return JsonInput{}, err
	}

	var jsonInput JsonInput
	err = json.Unmarshal([]byte(input), &jsonInput)
	if err != nil {
		return JsonInput{}, err
	}
	return jsonInput, nil
}

func (ji JsonInput) InputNames() ([]string, error) {
	convertedNames, err := convertSlice(ji.Names)
	if err != nil {
		return nil, err
	}
	return convertedNames, nil
}

func (ji JsonInput) InputTurns() (int, error) {
	return ji.Turns, nil
}
