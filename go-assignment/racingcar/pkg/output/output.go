package output

import (
	"errors"
	"fmt"
	"racing-car/racingcar/pkg"
	"sort"
	"strings"
)

const (
	printingNameLength = 5
)

type Output interface {
	PrintRank([]*pkg.User) error
}

var MaxRank int

func CreateOutput(ioType pkg.IoType) (Output, error) {
	switch {
	case pkg.Cli == ioType:
		return CliOutput{}, nil
	case pkg.File == ioType:
		return FileOutput{}, nil
	case pkg.Json == ioType:
		return JsonOutput{}, nil
	default:
		fmt.Println("err")
		return nil, errors.New("wrong ioType")
	}
}

func sortUsers(users []*pkg.User) {
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

func parseNameByLength(user *pkg.User) string {
	if len(user.Name) < printingNameLength+1 {
		return user.Name + strings.Repeat(" ", printingNameLength-len(user.Name))
	}
	return user.Name[:printingNameLength]
}

func parseWinners(sortedUsers []*pkg.User) ([]pkg.User, error) {
	topTurns := sortedUsers[0].NumberOfTurns
	winners := []pkg.User{*sortedUsers[0]}

	for i := 1; i < len(sortedUsers); i++ {
		switch {
		case sortedUsers[i].NumberOfTurns > topTurns:
			return nil, errors.New("not sorted")
		case topTurns == sortedUsers[i].NumberOfTurns:
			winners = append(winners, *sortedUsers[i])
		}
	}

	return winners, nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
