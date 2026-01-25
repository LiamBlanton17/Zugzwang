package engine

/*
This package contains all the aliases and definitions for the types and constant values used by the chess engine.
*/

// Aliasing FEN to string for better type saftey
type FEN string

// Aliasing PGN to string for better type safety
type PGN string

// Aliasing Piece to uint8 for better type safety
type Piece uint8

// Aliasing Color to unit8 for better type safety
type Color uint8

// Aliasing Eval to int32 for better type safety
type Eval int32

// Aliasing Move to uint32 for better type safety
type Move uint32

// Aliasing a Move, Eval pair
type MoveEval struct {
	move Move
	eval Eval
}

// Aliasing Square to unit8 for better type safety
// Uint8 has a max value of 256, enough to store all 64 possible squares in it
type Square uint8

// Aliasing a Zobrist hash to uint64 for better type safety
type ZobristHash uint8

// Aliasing game history to an array of Zobrist hashs
// This is the most effective way to use game history to check for repititions
type GameHistory []ZobristHash

// Aliasing BitBoard to uint64 for better type safety
// Bitboards are 64 bits of 0s meaning no piece and 1s meaning a piece
type BitBoard uint64

// Defining the number of pieces and colors
const (
	NUM_PIECES uint8 = 6
	NUM_COLORS uint8 = 2
)

// Defining the types of pieces
const (
	PAWN Piece = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

// Defining the types of colors
const (
	WHITE Color = iota
	BLACK
	EITHER_COLOR
)

// Defining the board structure
type Board struct {

	// Used to track where pieces, accessed by color and piece
	// Ex. Board.Pieces[WHITE][KNIGHT] is a bitboard where the 1s are the position of the White Knights
	Pieces [NUM_COLORS][NUM_PIECES]BitBoard

	// Used to track if the squares are occupied at all, by color or by either color
	// Ex. Board.Occupancy[EITHER_COLOR] is a bitboard where the 1s are occupied squares, by either color
	Occupancy [NUM_COLORS + 1]BitBoard

	// Check whose turn it is
	// Ex. Board.Turn == WHITE is true if it is whites turn, false if it is blacks
	Turn Color

	// Keep track of which square is the en passent square
	EPS Square

	// Keep track of the fifty move rule
	// This is incremented each time a move is played that is not a capture or a pawn move
	// If a capture or a pawn move is played, this counter is reset to 0
	// Unit8 is large enough as it will never be larger than 52
	HMC uint8

	// Keep track of the full move of the current position
	// This starts at 1 and is incremented after blacks move
	// The maximum number of moves is 8,848.5 under current FIDE rules, so unint16 is large enough
	FMC uint16

	// This is the Zobrist hash of the board position
	// This is vital in TT tables and hashing
	Zobrist ZobristHash
}
