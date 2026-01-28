package engine

import (
	"fmt"
	"math/bits"
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

// Converts a Square (a number that will map to 0-63 on a bitboard) to a string (such as h7 or d3)
func (sq Square) toString() string {
	// 255 is chosen to be the "NULL" value for the square
	if sq == NO_SQUARE {
		return " "
	}

	// Check the column
	str := ""
	switch sq % 8 {
	case 0:
		str = "a"
	case 1:
		str = "b"
	case 2:
		str = "c"
	case 3:
		str = "d"
	case 4:
		str = "e"
	case 5:
		str = "f"
	case 6:
		str = "g"
	case 7:
		str = "h"
	}

	// Check the row
	switch sq / 8 {
	case 0:
		str += "1"
	case 1:
		str += "2"
	case 2:
		str += "3"
	case 3:
		str += "4"
	case 4:
		str += "5"
	case 5:
		str += "6"
	case 6:
		str += "7"
	case 7:
		str += "8"
	}

	return str
}

// Converts a Piece to a character
func (p Piece) toString(c Color) string {
	switch p {
	case KING:
		if c == BLACK {
			return "k"
		}
		return "K"
	case QUEEN:
		if c == BLACK {
			return "q"
		}
		return "Q"
	case ROOK:
		if c == BLACK {
			return "r"
		}
		return "R"
	case BISHOP:
		if c == BLACK {
			return "b"
		}
		return "B"
	case KNIGHT:
		if c == BLACK {
			return "n"
		}
		return "N"
	case PAWN:
		if c == BLACK {
			return "p"
		}
		return "P"
	}

	return "-"
}

// Converts a Move to a string - allways sets the promotion piece to WHITE
func (m Move) toString() string {
	start := m.start.toString()
	target := m.target.toString()
	promotion := m.promotion.toString(WHITE)
	codeStr := "MOVE_CODE_NONE"
	switch m.code {
	case MOVE_CODE_CAPTURE:
		codeStr = "MOVE_CODE_CAPTURE"
	case MOVE_CODE_EN_PASSANT:
		codeStr = "MOVE_CODE_EN_PASSANT"
	case MOVE_CODE_DOUBLE_PAWN_PUSH:
		codeStr = "MOVE_CODE_DOUBLE_PAWN_PUSH"
	case MOVE_CODE_CASTLE:
		codeStr = "MOVE_CODE_CASTLE"
	}

	return fmt.Sprintf("(%v to %v) (promotion: %v) (code: %v)", start, target, promotion, codeStr)
}

// Helper function to get the index on a bitboard given some square number
func (s Square) bitBoardPosition() BitBoard {
	return BitBoard(uint64(1) << s)
}

// Helper function to get the LSB (Square) of a bitboard and pop it off
func (bb *BitBoard) popSquare() Square {
	// Index of the rightmost set bit (Least Significant Bit)
	idx := bits.TrailingZeros64(uint64(*bb))

	// Removes the lowest set bit
	*bb &= (*bb - 1)

	return Square(idx)
}

// Helper function to clear a bit off a bitboard
func (bb *BitBoard) clear(b BitBoard) {
	*bb &= ^b
}

// Helper function to set a bit on a bitboard
func (bb *BitBoard) set(b BitBoard) {
	*bb |= b
}

// Helper function to setup magic bitboards
func SetMaskOccupancy(index int, bitsInMask int, attackMask BitBoard) BitBoard {
	occupancy := BitBoard(0)

	for i := range bitsInMask {
		square := attackMask.popSquare()

		// If the i-th bit of 'index' is set, place a piece there
		if (index & (1 << i)) != 0 {
			occupancy |= (1 << square)
		}
	}

	return occupancy
}

// Globals and function to setup for Zobrist hashing
const MASTER_ZOBRIST = 20240928                                    // Used for initializing Zobrist values
var PIECE_ZOBRIST [NUM_COLORS][NUM_PIECES][NUM_SQUARES]ZobristHash // Global for Zobrist hashing pieces
var BLACK_TO_MOVE_ZOBRIST ZobristHash                              // Global for Zorbist hasing black to move
var CASTLING_ZOBRIST [16]ZobristHash                               // Global for Zobrist hashing castling rights (1 for each combination)
var ENPASSENT_ZOBRIST [8]ZobristHash                               // Global for Zobrist hashing enpassent column (8 columns totoal)
func initZobrist() {
	// Setup determistic hashing with one constant master key
	source := rand.NewSource(MASTER_ZOBRIST)
	rng := rand.New(source)

	// Setup piece hashing
	for color := range NUM_COLORS {
		for piece := range NUM_PIECES {
			for square := range NUM_SQUARES {
				PIECE_ZOBRIST[color][piece][square] = ZobristHash(rng.Uint64())
			}
		}
	}

	// Setup to move hashing
	BLACK_TO_MOVE_ZOBRIST = ZobristHash(rng.Uint64())

	// Setup castling hashing
	for i := range CASTLING_ZOBRIST {
		CASTLING_ZOBRIST[i] = ZobristHash(rng.Uint64())
	}

	// Setup en passent hashing
	for i := range ENPASSENT_ZOBRIST {
		ENPASSENT_ZOBRIST[i] = ZobristHash(rng.Uint64())
	}
}
