package pkg

import (
	"fmt"
	"strings"
)

const printingNameLength = 5

func PrintRank(users []*User) {
	topMinUsersNum := calculateMin(3, len(users))
	for _, user := range users[:topMinUsersNum] {
		fmt.Printf("%s:%s\n", parseNameByLength(user, printingNameLength), parseDashByNumberOfTurns(user))
	}
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
