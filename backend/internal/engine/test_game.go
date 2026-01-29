package engine

import (
	"cmp"
	"fmt"
	"slices"
)

/*
This file holds the function for testing the engine in a game where it just plays itself.
It prints the board, searches to a defined depth, plays that move, searches again.
Until mate.
*/

func TestGame() {
	fmt.Println("Starting test game.")
	fmt.Println()

	// Init the engine
	InitEngine()

	// Allocate the moveStack
	moveStack := make([][]Move, MAX_PLY)
	for i := range moveStack {
		moveStack[i] = make([]Move, 256)
	}

	// Setup the starting board
	board, _ := STARTING_POSITION_FEN.toBoard(nil)
	const depth = uint8(7)

	for {

		// Print board, search
		board.print()
		result := board.rootSearch(depth, moveStack, false)
		moveResults := result.moves

		// No moves possible, game over
		if len(moveResults) == 0 {
			fmt.Println("Game Over.")
			return
		}

		// Sort and get best result
		slices.SortFunc(moveResults, func(a, b MoveEval) int {
			return cmp.Compare(b.eval, a.eval)
		})

		// Eval needs to be context aware
		bestEval := moveResults[0].eval
		if board.Turn == BLACK {
			bestEval *= -1
		}
		bestMove := moveResults[0].move
		fmt.Printf("Best move: %v with eval %.3f\n", bestMove.toString(), float32(bestEval)/100)
		fmt.Printf("The engine searched: %d nodes\n\n\n", result.nodes)

		// Make the move
		board.makeMove(bestMove)

		break
	}
}
