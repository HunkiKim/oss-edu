package main

import (
	"log"
	"racing-car/racingcar/pkg"
)

func main() {
	names, err := pkg.InputNames()
	check(err)

	turns, err := pkg.InputTurns()
	check(err)

	users := pkg.CreateUsers(names)

	pkg.DoRace(users, turns)

	pkg.PrintRank(users)
}

func check(err error) {
	if err != nil {
		log.Fatalf("실패 %v", err)
	}
}
