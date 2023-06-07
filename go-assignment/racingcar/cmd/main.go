package main

import (
	"log"
	"racing-car/racingcar/pkg"
)

func main() {
	names, numberOfTurns, err := pkg.InputAll()
	if err != nil {
		log.Fatalf("실패 %v", err) // 다시확인
	}

	users := pkg.CountNumberOfTurns(names, numberOfTurns) // numberOfTurns 이름 바꾸기

	pkg.PrintRank(users)
}
