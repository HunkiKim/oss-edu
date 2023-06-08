package pkg

import (
	"fmt"
	"sort"
	"strings"
)

const printingNameLength = 5

func PrintRank(users []*User) {
	sortUsers(users)

	topMinUsersNum := calculateMin(3, len(users))
	for _, user := range users[:topMinUsersNum] {
		fmt.Printf("%s:%s\n", parseNameByLength(user, printingNameLength), parseDashByNumberOfTurns(user))
	}
}

func sortUsers(users []*User) {
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
func calculateMin(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

func parseNameByLength(user *User, length int) string {
	if len(user.name) < length+1 {
		return user.name + strings.Repeat(" ", length-len(user.name))
	}
	return user.name[:5]
}

func parseDashByNumberOfTurns(user *User) string {
	return strings.Repeat("-", user.numberOfTurns)
}
