package engine

import (
	"fmt"
	"time"
)

/*
This file holds the Preft test, which can be called via the command line.
This is used for validating the engine's correctness when it comes to move generation and searching
*/
func Preft() {
	fmt.Println("Starting PREFT test.")

	// The depth in ply to test too
	const TEST_DEPTH = 6

	// Init the engine
	InitEngine()

	// Allocate the moveStack
	moveStack := make([][]Move, TEST_DEPTH+1)
	for i := range moveStack {
		moveStack[i] = make([]Move, 256)
	}

	// Setup board from test position
	board, err := FEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0").toBoard(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	start := time.Now()
	result := board.negamax(TEST_DEPTH, moveStack)
	searchTime := time.Since(start)
	mnps := (float64(result.nodes) / searchTime.Seconds()) / 1000000

	fmt.Printf("Total nodes searched: %v\n", result.nodes)
	fmt.Printf("Total search time: %v\n", searchTime)
	fmt.Printf("Million nodes per second: %.3f\n", mnps)
}
