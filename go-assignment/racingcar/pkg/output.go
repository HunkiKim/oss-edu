package pkg

import "fmt"

func PrintRank(users []User) {
	for i := 0; i < calculateMin(3, len(users)); i++ { // 3이 아니라 나중에 입력받아서 처리하도록 변경 예정
		fmt.Printf("%s:%s", printName(&users[i]), printNumberOfTurns(users[i]))
	}
}

func calculateMin(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

func printName(user *User) string {
	if len(user.name) <= 5 {
		return user.name
	}
	return user.name[:5]
}

func printNumberOfTurns(user User) string {
	var answer string
	for i := 0; i < user.numberOfTurns; i++ {
		answer += "-"
	}
	return answer
}
