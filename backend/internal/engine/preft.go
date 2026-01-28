package engine

import "fmt"

/*
This file holds the Preft test, which can be called via the command line.
This is used for validating the engine's correctness when it comes to move generation and searching
*/
func Preft() {
	fmt.Println("Starting PREFT test.")

	// The depth in ply to test too
	const TEST_DEPTH = 5

	// Init the engine
	InitEngine()

	// Allocate the moveStack
	moveStack := make([][]Move, TEST_DEPTH+1)
	for i := range moveStack {
		moveStack[i] = make([]Move, 256)
	}

	// Setup board from starting position, do not worry about the error
	board, _ := STARTING_POSITION_FEN.toBoard(nil)

	board.print()
	result := board.negamax(TEST_DEPTH, moveStack)

	fmt.Printf("Total nodes searched: %v", result.nodes)
}
