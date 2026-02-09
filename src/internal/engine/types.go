package engine

/*
This package contains all the aliases and definitions for the types and constant values used by the chess engine.
*/

// Aliasing FEN to string for better type saftey
// The starting FEN position is: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
type FEN string

// Define the starting position
const STARTING_POSITION_FEN = FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

// Aliasing PGN to string for better type safety
type PGN string

// Aliasing Piece to uint8 for better type safety
type Piece uint8

// Aliasing Color to unit8 for better type safety
type Color uint8

// Aliasing Eval to int16 for better type safety
type Eval int16

// Definging the move structure
// Code will be a variety of things (is a caputure, white castle kingside, etc)
// promotion holds the piece to promote too
type Move struct {
	start     Square
	target    Square
	promotion Piece
	code      uint8
}

// This structure is used to unmake moves in place on a board, after making a move
type MoveUndo struct {
	cr          uint8
	hmc         uint8
	code        uint8
	isPromotion bool
	captured    Piece
	eps         Square
	start       Square
	target      Square
}

// Move code definitions
const (
	MOVE_CODE_NONE uint8 = iota
	MOVE_CODE_CAPTURE
	MOVE_CODE_EN_PASSANT
	MOVE_CODE_DOUBLE_PAWN_PUSH
	MOVE_CODE_CASTLE
)

// Aliasing a Move, Eval pair
type MoveEval struct {
	move Move
	eval Eval
}

// Aliasing killer moves for better type safety
type Killers [MAX_PLY][2]Move

// Aliasing cutoff history hueristic for better type safety
type CutoffHeuristic [NUM_COLORS][NUM_SQUARES][NUM_SQUARES]int

// Defining mins and maxes for the eval type, this is close to max for 16-bit int but not there (to avoid overflow issues)
const (
	MAX_EVAL = Eval(27000)
	MIN_EVAL = Eval(-27000)
)

// Define the max ply the engine will search too
const MAX_PLY = 64

// Defining the max number of moves in a position
// This comes from lichess official study that it is 218, but setting to 256 is fine
const MAX_NUMBER_OF_MOVES_IN_A_POSITION = 256

// Defining a delta margin to use in Delta Pruning in the Quiescence search
// This is in centipawns
const DELTA_MARGIN = 150

// Defining the starting history length
// This can be tweaked if needed, but shouldn't have too much of an effect on the performance
const STARTING_HISTORY_LENGTH = 50

// Defining the game stages
// These are used for helping the engine make more accruate evaluations of the position
const (
	OPENING uint8 = iota
	ENDGAME
)

// Aliasing Square to unit8 for better type safety
// Uint8 has a max value of 255, enough to store all 64 possible squares in it
// 255 is chosen to be the "NULL" value for the square
// Use square.bitBoardPosition to get the correct uint64 offset
type Square uint8

// Aliasing a Zobrist hash to uint64 for better type safety
type ZobristHash uint64

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
	NO_PIECE
)

// Defining the types of colors
const (
	WHITE Color = iota
	BLACK
	EITHER_COLOR
)

// Defining castling right constants
// This will help for readability when checking for castling rights
const (
	CASTLE_WK = 0x08
	CASTLE_WQ = 0x04
	CASTLE_BK = 0x02
	CASTLE_BQ = 0x01
)

// Defining characters of pieces, helpes with converting from strings into a board
const (
	CHAR_WK = 'K'
	CHAR_WQ = 'Q'
	CHAR_WR = 'R'
	CHAR_WB = 'B'
	CHAR_WN = 'N'
	CHAR_WP = 'P'
	CHAR_BK = 'k'
	CHAR_BQ = 'q'
	CHAR_BR = 'r'
	CHAR_BB = 'b'
	CHAR_BN = 'n'
	CHAR_BP = 'p'
)

// Defining the corners of the board as ints on a bit board for better readability
// Defining total number of squares on chess board (not strictly needed but can make some code look nicer)
// Defining some all squares here
const (
	NUM_SQUARES        = 64
	NO_SQUARE   Square = 255
)
const (
	A1 Square = 0
	B1 Square = 1
	C1 Square = 2
	D1 Square = 3
	E1 Square = 4
	F1 Square = 5
	G1 Square = 6
	H1 Square = 7

	A2 Square = 8
	B2 Square = 9
	C2 Square = 10
	D2 Square = 11
	E2 Square = 12
	F2 Square = 13
	G2 Square = 14
	H2 Square = 15

	A3 Square = 16
	B3 Square = 17
	C3 Square = 18
	D3 Square = 19
	E3 Square = 20
	F3 Square = 21
	G3 Square = 22
	H3 Square = 23

	A4 Square = 24
	B4 Square = 25
	C4 Square = 26
	D4 Square = 27
	E4 Square = 28
	F4 Square = 29
	G4 Square = 30
	H4 Square = 31

	A5 Square = 32
	B5 Square = 33
	C5 Square = 34
	D5 Square = 35
	E5 Square = 36
	F5 Square = 37
	G5 Square = 38
	H5 Square = 39

	A6 Square = 40
	B6 Square = 41
	C6 Square = 42
	D6 Square = 43
	E6 Square = 44
	F6 Square = 45
	G6 Square = 46
	H6 Square = 47

	A7 Square = 48
	B7 Square = 49
	C7 Square = 50
	D7 Square = 51
	E7 Square = 52
	F7 Square = 53
	G7 Square = 54
	H7 Square = 55

	A8 Square = 56
	B8 Square = 57
	C8 Square = 58
	D8 Square = 59
	E8 Square = 60
	F8 Square = 61
	G8 Square = 62
	H8 Square = 63
)

// Defining the board structure
type Board struct {

	// Used to track where pieces, accessed by color and piece
	// Ex. Board.Pieces[WHITE][KNIGHT] is a bitboard where the 1s are the position of the White Knights
	Pieces [NUM_COLORS][NUM_PIECES]BitBoard

	// Used to track if the squares are occupied at all, by color or by either color
	// Ex. Board.Occupancy[EITHER_COLOR] is a bitboard where the 1s are occupied squares, by either color
	Occupancy [NUM_COLORS + 1]BitBoard

	// Used to fast lookup of pieces, keeps track of what piece is on what square
	// This is a "mailbox" approach and trades a bit more space for more efficient piece look ups
	MailBox [NUM_SQUARES]Piece

	// Check whose turn it is
	// Ex. Board.Turn == WHITE is true if it is whites turn, false if it is blacks
	Turn Color

	// Keep track of castling rights
	// 0000 KQkq
	// lsb is black queenside, then black kingside, white queen, white kingside
	CR uint8

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

	// History of board positions
	// This is allocated once at the start of the search
	// These are just Zobrist hashes
	History GameHistory

	// This stores the square of the king for both sides
	// Keep this updated, makes finding the king more efficient during move generation
	KingSquare [NUM_COLORS]Square
}
