package engine

import (
	"fmt"
	"math/rand"
)

/*
This file contains random helper functions that don't belong anywhere else
*/

// Converts a string (such as h7 or d3) to a Square (a number that will map to 0-63 on a bitboard)
func stringToSquare(s string) (Square, error) {
	// 255 is chosen to be the "NULL" value for the square
	if s == "-" || s == "" || s == " " {
		return Square(255), nil
	}

	// Once the "NULL" value is handled, square strings must have a length of 2
	// First is the column (a, b, c, etc) second is the row (1, 2, 3, etc)
	if len(s) != 2 {
		return Square(255), fmt.Errorf("Invalid square length of %v (should be 2 [col:row])", len(s))
	}

	// Check the column
	col := 0
	switch s[0] {
	case 'a':
		col = 0
	case 'b':
		col = 1
	case 'c':
		col = 2
	case 'd':
		col = 3
	case 'e':
		col = 4
	case 'f':
		col = 5
	case 'g':
		col = 6
	case 'h':
		col = 7
	default:
		return Square(255), fmt.Errorf("Invalid square column: %v", s[0])
	}

	// Check the row
	// This is done verbosely instead of converting string to int to check for more errors
	row := 0
	switch s[1] {
	case '1':
		row = 0
	case '2':
		row = 1
	case '3':
		row = 2
	case '4':
		row = 3
	case '5':
		row = 4
	case '6':
		row = 5
	case '7':
		row = 6
	case '8':
		row = 7
	default:
		return Square(255), fmt.Errorf("Invalid square row: %v", s[1])
	}

	// Combine row and column into one square
	return Square(row*8 + col), nil
}

// Helper function to get the index on a bitboard given some square number
func (s Square) bitBoardPosition() BitBoard {
	return BitBoard(uint64(1) << s)
}

// Globals and function to setup for Zobrist hashing
const MASTER_ZOBRIST = 20240928                                 // Used for initializing Zobrist values
var PIECE_ZOBRIST [NUM_COLORS][NUM_PIECES][TOTAL_SQUARES]uint64 // Global for Zobrist hashing pieces
var TO_MOVE_ZOBRIST uint64                                      // Global for Zorbist hasing black to move
var CASTLING_ZOBRIST [16]uint64                                 // Global for Zobrist hashing castling rights (1 for each combination)
var ENPASSENT_ZOBRIST [8]uint64                                 // Global for Zobrist hashing enpassent column (8 columns totoal)
func initZobrist() {
	// Setup determistic hashing with one constant master key
	source := rand.NewSource(MASTER_ZOBRIST)
	rng := rand.New(source)

	// Setup piece hashing
	for color := range NUM_COLORS {
		for piece := range NUM_PIECES {
			for square := range TOTAL_SQUARES {
				PIECE_ZOBRIST[color][piece][square] = rng.Uint64()
			}
		}
	}

	// Setup to move hashing
	TO_MOVE_ZOBRIST = rng.Uint64()

	// Setup castling hashing
	for i := range CASTLING_ZOBRIST {
		CASTLING_ZOBRIST[i] = rng.Uint64()
	}

	// Setup en passent hashing
	for i := range ENPASSENT_ZOBRIST {
		ENPASSENT_ZOBRIST[i] = rng.Uint64()
	}
}
