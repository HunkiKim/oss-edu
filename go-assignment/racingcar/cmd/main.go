package main

import (
	"errors"
	"log"
	"sort"

	"github.com/spf13/cobra"

	"racing-car/pkg/command"
	"racing-car/pkg/file"
	"racing-car/pkg/interfaces"
	"racing-car/pkg/json"
	"racing-car/pkg/user"
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
			cmd.SilenceUsage = false
			return startApplication(flags)
		},
	}

	rootCmd.Flags().StringVar(&flags.format, "format", "", "")
	rootCmd.Flags().IntVarP(&flags.maxRank, "maxRank", "", 3, "")

	return rootCmd
}

func startApplication(f *racingFlags) error {
	reader, err := NewReader(f.format)
	if err != nil {
		return err
	}

	names, turns, err := reader.Read()
	if err != nil {
		return err
	}

	users, err := user.CreateUsers(names)
	if err != nil {
		return err
	}

	user.DoRace(users, turns)

	writer, err := NewWriter(f.format, "./result.txt")
	if err != nil {
		return err
	}

	sortUsers(users)

	winners, err := parseWinners(users)
	if err != nil {
		return err
	}
	topUsers := users[:min(f.maxRank, len(users))]
	if err = writer.Write(winners, topUsers); err != nil {
		return err
	}

	return nil
}

func sortUsers(users []*user.User) {
	sort.Slice(users, func(i, j int) bool {
		switch {
		case users[i].NumberOfTurns > users[j].NumberOfTurns:
			return true
		case users[i].NumberOfTurns < users[j].NumberOfTurns:
			return false
		case users[i].Name < users[j].Name:
			return true
		default:
			return false
		}
	})
}

func parseWinners(sortedUsers []*user.User) ([]*user.User, error) {
	topTurns := sortedUsers[0].NumberOfTurns
	cnt := 1

	for _, u := range sortedUsers {
		if u.NumberOfTurns != topTurns {
			break
		}
		cnt += 1
	}

	return sortedUsers[:cnt], nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func main() {
	var cmd = InitRacingCmd()
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func NewReader(format string) (interfaces.Reader, error) {
	switch format {
	case "Command":
		return command.NewReader(), nil
	case "File":
		return file.NewReader(), nil
	case "Json":
		return json.NewReader(), nil
	default:
		return nil, errors.New("wrong format")
	}
}

func NewWriter(format, path string) (interfaces.Writer, error) {
	switch format {
	case "Command":
		return command.NewWriter(), nil
	case "File":
		return file.NewWriter(path), nil // filepath, 환경변수 외부에서 주입 , 캡슐화 하는 방법 중 하나, 추상화
	case "Json":
		return json.NewWriter(), nil
	default:
		return nil, errors.New("wrong format")
	}
}
