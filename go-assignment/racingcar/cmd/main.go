package main

import (
	"github.com/spf13/cobra"
	"log"
	"sort"

	"racing-car/racingcar/pkg"
	"racing-car/racingcar/pkg/user"
	"racing-car/racingcar/pkg/util"
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
	reader, err := util.NewReader(f.format)
	if err != nil {
		return err
	}

	names, turns, err := reader.Read()
	if err != nil {
		return err
	}

	users, err := pkg.CreateUsers(names)
	if err != nil {
		return err
	}

	pkg.DoRace(users, turns)

	writer, err := util.NewWriter(f.format)
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
