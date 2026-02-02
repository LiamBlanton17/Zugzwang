package engine

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"time"
)

/*
This file holds the testing for the strength of the engine.
It holds the tests to be run and see the results, to know if changes improve the engine or not.
*/

type StrengthTestTest struct {
	fen           FEN
	stockfishEval Eval
	stockfishMove string
	depth         uint8
	rounds        int
}

func StrengthTest() {
	fmt.Println("Starting strength test.")
	fmt.Println()

	// Init the engine
	InitEngine()

	// Positions to run the strength test on.
	var positions []StrengthTestTest = []StrengthTestTest{
		{
			fen:           STARTING_POSITION_FEN,
			stockfishEval: 20,
			stockfishMove: "c4",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "r1b1k2r/pp1n2pp/1qn1pp2/3pP3/1b1P1P2/3B1N2/PP1B2PP/R2QK1NR w KQkq - 4 11",
			stockfishEval: 169,
			stockfishMove: "Ne2",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "r1b3k1/pp1nb1pp/1q2p3/3pP3/3n4/P2B1P2/1PQBN2P/R3K2R w KQ-- - 0 16",
			stockfishEval: 0,
			stockfishMove: "Nxd4",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "r1b4k/pp4pB/4pB2/3p4/2n2P1q/P7/1PQ4P/1K1R3R b ---- - 0 22",
			stockfishEval: 176,
			stockfishMove: "Qxh7",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "r3k2r/1b3ppp/pq2p3/1pb5/P5n1/3B1N2/1PP1QPPP/R1B2RK1 b --kq - 6 16",
			stockfishEval: -38,
			stockfishMove: "b4",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "2rr4/1b2kppp/p3p3/P1n1N3/1pB5/1P2P2P/2P3P1/R2R2K1 b ---- - 0 27",
			stockfishEval: -53,
			stockfishMove: "be4",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "8/2k4p/B1b1p1p1/5pP1/7R/1P2P2P/2r5/4K3 w ---- - 0 40",
			stockfishEval: 7,
			stockfishMove: "Rxh7",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "rnbqkb1r/pp2pp1p/5np1/3P4/8/2N5/PP1P1PPP/R1BQKBNR w KQkq - 0 6",
			stockfishEval: 118,
			stockfishMove: "Bc4",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "r2q1rk1/1p2ppb1/6pp/p1nP1b2/P1PN4/1QN1B3/1P3PPP/R3R1K1 w ---- - 1 16",
			stockfishEval: 332,
			stockfishMove: "Qb5",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "2rqr1k1/p4p2/1p2p1p1/4Nn2/3PR3/P1PQ4/5PP1/2R3K1 b ---- - 0 23",
			stockfishEval: -93,
			stockfishMove: "Kg7",
			depth:         7,
			rounds:        3,
		},
		{
			fen:           "2rr4/p4pk1/1p2p1pn/4N3/3P4/P1PR4/5PP1/3R1K2 b ---- - 2 30",
			stockfishEval: -92,
			stockfishMove: "g5",
			depth:         7,
			rounds:        3,
		},
	}

	// Setup profiler
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	for pi, position := range positions {
		// Setup the starting board
		fmt.Printf("Starting test of position %d.\n", pi+1)
		board, _ := position.fen.toBoard(nil)
		depth := position.depth
		aggSearchTime := int64(0)
		nodes := 0
		bestEval := Eval(0)
		rounds := position.rounds
		bestMove := Move{}
		for rounds > 0 {

			// Clear TT (otherwise the entire search gets cached essentially)
			ClearTT()

			// search
			timeStart := time.Now()
			result := board.rootSearch(depth, false)
			moveResults := result.moves
			aggSearchTime += time.Since(timeStart).Milliseconds()
			nodes = result.nodes

			// Sort and get best result
			slices.SortFunc(moveResults, func(a, b MoveEval) int {
				return cmp.Compare(b.eval, a.eval)
			})

			// Eval needs to be context aware
			bestEval = moveResults[0].eval
			if board.Turn == BLACK {
				bestEval *= -1
			}
			bestMove = moveResults[0].move

			rounds--
		}

		// Print final results
		aggSearchTime /= int64(position.rounds)
		nps := float64(nodes) / (float64(aggSearchTime) / 1000.0)
		mnps := nps / 1_000_000.0
		board.print()
		fmt.Printf("Stockfish move %v: eval %.3f\n", position.stockfishMove, float32(position.stockfishEval)/100)
		fmt.Printf("Zugzwang move %v: eval %.3f\n", bestMove.toString(), float32(bestEval)/100)
		fmt.Printf("The engine searched: %d nodes\n", nodes)
		fmt.Printf("The time searched was: %d milliseconds\n", aggSearchTime)
		fmt.Printf("The Mn/s was: %.3f\n\n", mnps)
	}
}
