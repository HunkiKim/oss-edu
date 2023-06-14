package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"racing-car/racingcar/pkg"
	"racing-car/racingcar/pkg/input"
	"racing-car/racingcar/pkg/output"
)

type racingFlags struct {
	format  string
	maxRank int
}

func InitRacingCmd() *cobra.Command {
	flags := &racingFlags{}
	var rootCmd = &cobra.Command{
		Use:   "racing start",
		Short: "Racing Car Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return startRacing(flags)
		},
	}

	rootCmd.Flags().StringVar(&flags.format, "format", "", "")
	rootCmd.Flags().IntVarP(&flags.maxRank, "maxRank", "", 3, "")

	return rootCmd
}

func startRacing(f *racingFlags) error {
	var ioType pkg.IoType
	switch f.format {
	case "Command":
		ioType = pkg.Cli
	case "File":
		ioType = pkg.File
	case "Json":
		ioType = pkg.Json
	default:
		return errors.New("잘못된 파일 형식")
	}

	createdInput, err := input.CreateInput(ioType)
	if err != nil {
		return fmt.Errorf("입출력 선택 에러: %v", err)
	}

	names, err := createdInput.InputNames()
	if err != nil {
		return fmt.Errorf("이름 입력 에러: %v", err)
	}

	turns, err := createdInput.InputTurns()
	if err != nil {
		return fmt.Errorf("도는 횟수 입력 에러: %v", err)
	}

	users := pkg.CreateUsers(names)

	pkg.DoRace(users, turns)

	createOutput, err := output.CreateOutput(ioType)
	if err != nil {
		return fmt.Errorf("output 생성 에러: %v", err)
	}

	output.MaxRank = f.maxRank
	err = createOutput.PrintRank(users)
	if err != nil {
		return fmt.Errorf("출력 에러: %v", err)
	}

	return nil
}

func main() {
	var cmd = InitRacingCmd()
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
