package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"racing-car/racingcar/pkg/interfaces"
)

type Reader struct {
	*bufio.Reader
}

func NewReader() interfaces.Reader {
	return &Reader{bufio.NewReader(os.Stdin)}
}

type jsonInput struct {
	Names string
	Turns int
}

func (r *Reader) Read() (string, int, error) {
	fmt.Print("입력:")
	bytes, err := r.inputJsonFormat()
	if err != nil {
		return "", -1, err
	}

	j := jsonInput{}
	err = json.Unmarshal([]byte(bytes), &j)
	if err != nil {
		return "", -1, err
	}

	return j.Names, j.Turns, nil
}

func (r *Reader) inputJsonFormat() ([]byte, error) {
	var input []byte

	input, _, err := r.ReadLine()
	if err != nil {
		fmt.Println("입력 에러:", err)
		return nil, err
	}

	return input, nil
}
