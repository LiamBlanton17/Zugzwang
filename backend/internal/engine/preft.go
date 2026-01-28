package engine

import (
	"fmt"
	"time"
)

/*
This file holds the Perft test, which can be called via the command line.
This is used for validating the engine's correctness when it comes to move generation and searching.
This function does not need tests, as its the test itself.
All tests are taken from: https://www.chessprogramming.org/Perft_Results
*/
type PerftTest struct {
	name          string
	position      FEN
	depth         uint8
	expectedNodes int
}

func Perft() {
	fmt.Println("Starting PREFT test.")

	// The Perft tests to run
	tests := []PerftTest{
		{
			name:          "Position 1: Starting Position", // Position 1 on the website
			position:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:         6,
			expectedNodes: 119060324,
		},
		{
			name:          "Position 2: Kiwipete", // Position 2 on the website
			position:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:         5,
			expectedNodes: 193690690,
		},
		{
			name:          "Position 3: Rook/Pawn Endgame", // Position 3 on the website
			position:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:         7,
			expectedNodes: 178633661,
		},
		{
			name:          "Position 4: Tricky Bug Catcher", // Position 5 on the website
			position:      "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ-- - 1 8",
			depth:         5,
			expectedNodes: 89941194,
		},
		{
			name:          "Position 5: Alternative Perft by Steven Edwards", // Position 6 on the website
			position:      "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			depth:         5,
			expectedNodes: 164075551,
		},
	}

	// Init the engine
	InitEngine()

	for _, test := range tests {
		// Allocate the moveStack
		moveStack := make([][]Move, MAX_PLY)
		for i := range moveStack {
			moveStack[i] = make([]Move, 256)
		}

		// Setup board from test position
		board, err := test.position.toBoard(nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Search
		start := time.Now()
		result := board.rootSearch(test.depth, moveStack, false)
		searchTime := time.Since(start)
		mnps := (float64(result.nodes) / searchTime.Seconds()) / 1000000

		// Print results
		fmt.Printf("%v\n", test.name)
		fmt.Printf("Total search time: %v\n", searchTime)
		fmt.Printf("Million nodes per second: %.3f\n", mnps)
		fmt.Printf("Total nodes searched: %v\n", result.nodes)
		fmt.Printf("Expected nodes: %v\n", test.expectedNodes)
		fmt.Printf("Passed test: %v\n\n\n", test.expectedNodes == result.nodes)
	}
}
