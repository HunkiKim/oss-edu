package main

import (
	"log"
	"racing-car/racingcar/pkg"
)

func main() {
	names, err := pkg.InputNames()
	if err != nil {
		log.Fatalf("이름 입력 에러: %v", err)
	}

	turns, err := pkg.InputTurns()
	if err != nil {
		log.Fatalf("도는 횟수 입력 에러: %v", err)
	}

	users := pkg.CreateUsers(names)

	pkg.DoRace(users, turns)

	pkg.PrintRank(users)
}
