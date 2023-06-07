package pkg

import "fmt"

func PrintRank(users []User) {
	for i := 0; i < calculateMin(3, len(users)); i++ { // 3이 아니라 나중에 입력받아서 처리하도록
		fmt.Println(printName(&users[i]) + printNumberOfTurns(users[i]))
		//fmt.Sprintf() // 이걸로 수정한번해보기
		// 리플렉션 최대한 쓰지 말기
	}
}

func calculateMin(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

func printName(user *User) string { // 복사 비용 줄이려면 포인터로 주소 전달하기
	if len(user.name) <= 5 {
		return user.name + ":"
	}
	return user.name[:5] + ":"
}

func printNumberOfTurns(user User) string {
	var answer string
	for i := 0; i < user.numberOfTurns; i++ {
		answer += "-"
	}
	return answer
}
