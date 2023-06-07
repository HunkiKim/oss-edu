package main

import (
	"log"
	"racing-car/racingcar/pkg"
)

func main() {
	names, turns, err := pkg.InputAll()
	if err != nil {
		log.Fatalf("실패 %v", err)
	}

	users := pkg.DoRace(names, turns)

	pkg.PrintRank(users)
}
