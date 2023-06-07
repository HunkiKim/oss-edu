package main

import (
	"log"
	"racing-car/racingcar/pkg"
)

func main() {
	userNames, turns, err := pkg.InputAll()
	if err != nil {
		log.Fatalf("실패 %v", err)
	}

	users := pkg.CreateUsers(userNames)

	pkg.DoRace(users, turns)

	pkg.SortUsers(users)

	pkg.PrintRank(users)
}
