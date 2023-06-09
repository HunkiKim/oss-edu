package pkg

import (
	"fmt"
	"sort"
	"strings"
)

const (
	printingNameLength = 5
	maxTopUsersNum     = 3
)

func PrintRank(users []*user) {
	sortUsers(users)

	var min = func(x, y int) int {
		if x > y {
			return y
		}
		return x
	}
	topUsers := users[:min(maxTopUsersNum, len(users))]
	for _, user := range topUsers {
		fmt.Printf("%s:%s\n", parseNameByLength(user), strings.Repeat("-", user.numberOfTurns))
	}
}

func sortUsers(users []*user) {
	sort.Slice(users, func(i, j int) bool {
		switch {
		case users[i].numberOfTurns > users[j].numberOfTurns:
			return true
		case users[i].numberOfTurns < users[j].numberOfTurns:
			return false
		case users[i].name < users[j].name:
			return true
		default:
			return false
		}
	})
}

func parseNameByLength(user *user) string {
	if len(user.name) < printingNameLength+1 {
		return user.name + strings.Repeat(" ", printingNameLength-len(user.name))
	}
	return user.name[:printingNameLength]
}
