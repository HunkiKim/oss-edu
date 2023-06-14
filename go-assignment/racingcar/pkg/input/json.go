package input

type JsonInput struct {
	Names string
	Turns int
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
