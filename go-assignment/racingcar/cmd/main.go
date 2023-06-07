package main

import (
	"log"
	"racing-car/racingcar/pkg"
)

func main() {
	userNames, turns, err := pkg.InputAll()
	if err != nil {
		log.Fatalf("실패 %v", err) // 다시확인
	}

	users := pkg.CreateUsers(userNames)

	racedUsers := pkg.DoRace(users, turns)

	pkg.PrintRank(racedUsers)
}
