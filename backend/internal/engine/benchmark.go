package engine

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

/*
This file handles the engine benchmarking
*/

type benchmarkTestCandidate struct {
	move   string
	points int
}

type benchmarkTest struct {
	fen          FEN
	id           string
	depth        uint8
	candidates   []benchmarkTestCandidate
	candidateStr string
}

// Called to load run the benchmark
func RunBenchmark() {
	fmt.Println("Starting the benchmark test.")

	// Init the engine
	InitEngine()

	// Load the tests
	tests := make([]benchmarkTest, 0, 1000)

	// FROM STS1 - Undermining
	f, err := os.Open("benchmarktests/STS1.epd")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tests = append(tests, loadTest(scanner.Text(), false))
	}

	// FROM STS7 - Offer of Simplification
	f, err = os.Open("benchmarktests/STS7.epd")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		tests = append(tests, loadTest(scanner.Text(), false))
	}

	// FROM STS6 - Recapturing
	f, err = os.Open("benchmarktests/STS6.epd")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		tests = append(tests, loadTest(scanner.Text(), false))
	}

	// FROM STS2 - Open Files and Diagonals (need to swap id and candidates)
	f, err = os.Open("benchmarktests/STS2.epd")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		tests = append(tests, loadTest(scanner.Text(), true))
	}

	// Strength test totals
	totalNodes := 0
	totalSearchTime := 0
	totalPoints := 0
	foundBestMove := 0
	foundACandidate := 0

	// Run the tests
	for _, test := range tests {
		// Clear TT (otherwise the bad results across tests)
		ClearTT()

		// Setup the starting board
		fmt.Printf("Test: %v\n", test.id)
		fmt.Printf("Candidates: %v\n", test.candidateStr)
		board, err := test.fen.toBoard(nil)
		if err != nil {
			fmt.Println(test.fen)
			panic(err)
		}
		board.print()

		depth := test.depth
		aggSearchTime := int64(0)
		nodes := 0
		bestEval := Eval(0)
		bestMove := Move{}

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

		// Get points
		points := 0
		bestMovePCN := bestMove.toPCN()
		for _, candidate := range test.candidates {
			if candidate.move == bestMovePCN {
				points = candidate.points
				foundACandidate++
				if points == 10 {
					foundBestMove++
				}
			}
		}

		// Print position results
		nps := float64(nodes) / (float64(aggSearchTime) / 1000.0)
		mnps := nps / 1_000_000.0
		fmt.Printf("Zugzwang move %v: eval %.3f\n", bestMovePCN, float32(bestEval)/100)
		fmt.Printf("Zugzwang points: %d\n", points)
		fmt.Printf("The engine searched: %d nodes\n", nodes)
		fmt.Printf("The time searched was: %d milliseconds\n", aggSearchTime)
		fmt.Printf("The Mn/s was: %.3f\n\n", mnps)

		// Update totals
		totalNodes += nodes
		totalSearchTime += int(aggSearchTime)
		totalPoints += points
	}

	// Print final results
	avgTotalPoints := float32(totalPoints) / float32(len(tests))
	avgTotalNodes := totalNodes / len(tests)
	avgTotalSearchTime := totalSearchTime / len(tests)
	avgTimesFoundCandidate := float32(foundACandidate) / float32(len(tests))
	avgTimesFoundBest := float32(foundBestMove) / float32(len(tests))
	nps := float64(totalNodes) / (float64(totalSearchTime) / 1000.0)
	mnps := nps / 1_000_000.0
	fmt.Printf("---------------------\nFinal Results\n---------------------\n")
	fmt.Printf("Average Points: %.2f\n", avgTotalPoints)
	fmt.Printf("Percentage found a candidate: %.2f\n", avgTimesFoundCandidate*100)
	fmt.Printf("Percentage found the best (top 1): %.2f\n", avgTimesFoundBest*100)
	fmt.Printf("Average Nodes: %d\n", avgTotalNodes)
	fmt.Printf("Average Search Time: %d\n", avgTotalSearchTime)
	fmt.Printf("Average Mn/s: %.2f\n\n", mnps)

}

// Called to load a single test line
func loadTest(line string, swapIdAndCandidates bool) benchmarkTest {
	parts := strings.Split(line, ";")
	if len(parts) != 4 {
		fmt.Println(line)
		panic("Failed to load a test.")
	}

	// Parse out parts
	almostFen := parts[0]
	idStr := parts[1]
	candidateStr := strings.Split(parts[2], "\"")[1]

	// This is for STS2.epd, where the file has id and candidates swapped lol idk why bro did that
	if swapIdAndCandidates {
		idStr = parts[2]
		candidateStr = strings.Split(parts[1], "\"")[1]
	}

	// Drop final two parts in fen, replaec with move counters (0 0)
	fen := FEN(strings.Join(strings.Split(almostFen, " ")[:4], " ") + " 0 0")
	board, err := fen.toBoard(nil)
	if err != nil {
		fmt.Println(fen)
		panic(err)
	}

	// Only keep id name (between parens)
	id := strings.Split(idStr, "\"")[1]

	// Split out candidates
	candidates := make([]benchmarkTestCandidate, 0, 10)
	for candidate := range strings.SplitSeq(candidateStr, ",") {
		parts := strings.Split(candidate, "=")

		if len(parts) == 2 {
			move := strings.TrimSpace(parts[0]) // correctly gets "f5"
			points := 0

			// Parse the integer from the second part
			fmt.Sscanf(parts[1], "%d", &points)

			// Now it is safe to convert
			pcn, err := board.SanToPCN(move)
			if err != nil {
				board.print()
				panic(err)
			}

			candidates = append(candidates, benchmarkTestCandidate{move: pcn, points: points})
		}
	}

	return benchmarkTest{
		fen:          fen,
		id:           id,
		candidates:   candidates,
		candidateStr: candidateStr,
		depth:        7,
	}
}
