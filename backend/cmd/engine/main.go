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
	flag.StringVar(&action, "action", "preft", "the action the program takes")

	switch action {
	case "preft":
		engine.Preft()
	default:
		fmt.Println("The action is not supported: ", action)
	}
}
