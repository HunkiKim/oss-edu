package pkg

import (
	"fmt"
	"strings"
)

func PrintRank(users []*User) {
	topMinUsersNum := calculateMin(3, len(users))
	for i := 0; i < topMinUsersNum; i++ { // 3이 아니라 나중에 입력받아서 처리하도록 변경 예정
		fmt.Printf("%s:%s\n", parseFiveLengthName(users[i]), parseDashByNumberOfTurns(users[i]))
	}
}

func calculateMin(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

func parseFiveLengthName(user *User) string {
	if len(user.name) <= 5 {
		return user.name + strings.Repeat(" ", 5-len(user.name))
	}
	return user.name[:5]
}

func parseDashByNumberOfTurns(user *User) string {
	return strings.Repeat("-", user.numberOfTurns)
}
