package main

import (
	"backend/internal/engine"
	"flag"
	"fmt"
)

/*
This is the main binary for the command line interface with the chess engine
*/

func main() {

	var action string
	flag.StringVar(&action, "action", "perft", "the action the program takes")
	flag.Parse()

	switch action {
	case "perft":
		engine.Perft()
	case "testgame":
		engine.TestGame()
	default:
		fmt.Println("The action is not supported: ", action)
	}
}
