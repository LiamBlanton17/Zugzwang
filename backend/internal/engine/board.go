package engine

import (
	"fmt"
	"strconv"
	"strings"
)

/*
This file contains functionality related to setup, searching and evalution of a board.
*/

// Take a FEN string and turn it into a board, ready for the engine to search over it
func (position FEN) toBoard() (*Board, error) {
	board := Board{}

	// The starting FEN position is: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	// Parts is split up:
	// 0: the pieces on the board
	// 1: the turn to move ('w' or 'b')
	// 2: the castling rights 'KQkq'
	// 3: the en passent square
	// 4: the half move clock
	// 5: the full move number
	parts := strings.Split(string(position), " ")
	if len(parts) != 6 {
		return nil, fmt.Errorf("Invalid FEN string; Should be 6 parts in a FEN string (found %v)", len(parts))
	}
	pieces, turn, castling, enpassent, halfMove, fullMore := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]

	// Setup pieces
	err := board.setPieces(pieces)
	if err != nil {
		return nil, err
	}

	// Setup turn
	switch turn {
	case "w":
		board.Turn = WHITE
	case "b":
		board.Turn = BLACK
	default:
		return nil, fmt.Errorf("Invalid FEN string; Turn is not white or black")
	}

	// Setup castling
	if len(castling) != 4 {
		return nil, fmt.Errorf("Invalid FEN string; Castling rights string is not of length 4")
	}
	// This is more condensed than 4 if statements, but it could accept invalid FENs
	// KQkq is the correct order, but this would accept them out of order, like kQKq
	// This is mostly fine. Maybe in future can just be verbose and check each character one by one
	for s, v := range map[rune]uint8{CHAR_WK: CASTLE_WK, CHAR_WQ: CASTLE_WQ, CHAR_BK: CASTLE_BK, CHAR_BQ: CASTLE_BQ} {
		if strings.Contains(castling, string(s)) {
			board.CR |= v
		}
	}

	// Setup en passent
	eps, err := stringToSquare(enpassent)
	if err != nil {
		return nil, err
	}
	board.EPS = eps

	// Setup halfMove
	hm, err := strconv.Atoi(halfMove)
	if err != nil {
		return nil, err
	}
	board.HMC = uint8(hm)

	// Setup full move
	fm, err := strconv.Atoi(fullMore)
	if err != nil {
		return nil, err
	}
	board.FMC = uint16(fm)

	// Setup the Zobrist hash
	board.Zobrist = board.toZobrist()

	return &board, nil
}

// Return the Zobrist hash of a board
func (b *Board) toZobrist() ZobristHash {
	var hash ZobristHash

	// Hash in the pieces
	for color := range NUM_COLORS {
		for piece := range NUM_PIECES {
			pieceBits := b.Pieces[color][piece]
			for pieceBits > 0 {
				sq := pieceBits.popSquare()
				hash ^= PIECE_ZOBRIST[color][piece][sq]
			}
		}
	}

	// Hash in the move
	if b.Turn == BLACK {
		hash ^= BLACK_TO_MOVE_ZOBRIST
	}

	// Hash in the castling rights
	hash ^= CASTLING_ZOBRIST[b.CR]

	// Hash in the en passent column
	// Only if EPS is present
	if b.EPS != NO_SQUARE {
		hash ^= ENPASSENT_ZOBRIST[b.EPS%8]
	}

	return hash
}

// Take a pieces string from a FEN notation and set the board's pieces/occupancy/etc
func (b *Board) setPieces(pieces string) error {
	// Split the string first and check its length
	pieceParts := strings.Split(pieces, "/")
	if len(pieceParts) != 8 {
		return fmt.Errorf("Invalid FEN string; Should be 8 rows of pieces in the FEN string (found %v)", len(pieceParts))
	}

	// Loop over each part and each character in the part
	// Check each character, and update the correct bit boards and advance the index as needed
	idx := Square(A8) // The current index we are at (starting in top left of board at 56)
	for _, part := range pieceParts {
		holdIdx := idx // Save this to drop it down a row, and for FEN validation
		for _, c := range part {
			// If err is not nil, then c was not an integer
			if c >= '1' && c <= '8' {
				idx += Square(c - '0')
			} else {
				position := idx.bitBoardPosition()
				switch c {
				case CHAR_BK:
					b.Pieces[BLACK][KING] |= position
					b.Occupancy[BLACK] |= position
				case CHAR_BQ:
					b.Pieces[BLACK][QUEEN] |= position
					b.Occupancy[BLACK] |= position
				case CHAR_BR:
					b.Pieces[BLACK][ROOK] |= position
					b.Occupancy[BLACK] |= position
				case CHAR_BB:
					b.Pieces[BLACK][BISHOP] |= position
					b.Occupancy[BLACK] |= position
				case CHAR_BN:
					b.Pieces[BLACK][KNIGHT] |= position
					b.Occupancy[BLACK] |= position
				case CHAR_BP:
					b.Pieces[BLACK][PAWN] |= position
					b.Occupancy[BLACK] |= position
				case CHAR_WK:
					b.Pieces[WHITE][KING] |= position
					b.Occupancy[WHITE] |= position
				case CHAR_WQ:
					b.Pieces[WHITE][QUEEN] |= position
					b.Occupancy[WHITE] |= position
				case CHAR_WR:
					b.Pieces[WHITE][ROOK] |= position
					b.Occupancy[WHITE] |= position
				case CHAR_WB:
					b.Pieces[WHITE][BISHOP] |= position
					b.Occupancy[WHITE] |= position
				case CHAR_WN:
					b.Pieces[WHITE][KNIGHT] |= position
					b.Occupancy[WHITE] |= position
				case CHAR_WP:
					b.Pieces[WHITE][PAWN] |= position
					b.Occupancy[WHITE] |= position
				default:
					return fmt.Errorf("Invalid FEN string; Invalid piece present on the board: %q", c)
				}
				b.Occupancy[EITHER_COLOR] |= position
				idx += 1
			}
		}

		// Make sure the index has moved 8 squares exactly, else invalidate the FEN string
		if idx-8 != holdIdx {
			return fmt.Errorf("Invalid FEN string; String did not move 8 columns in a row (moved %v)", idx-holdIdx)
		}

		// Now dropping down a row (8 moves it down 1 whole row)
		idx = holdIdx - 8
	}

	return nil
}

// Helper function to get the pieces of the opposing player of the current turn
func (b *Board) getEnemyPieces() BitBoard {
	if b.Turn == WHITE {
		return b.Occupancy[BLACK]
	}
	return b.Occupancy[WHITE]
}

func buildGameHistory(history []FEN) (*GameHistory, error) {
	var gameHistory GameHistory

	for _, h := range history {
		b, err := h.toBoard()
		if err != nil {
			return nil, err
		}

		gameHistory = append(gameHistory, b.toZobrist())
	}

	return &gameHistory, nil
}

type BoardSearchResults struct {
	Nodes     int32
	MoveEvals []MoveEval
}

func (b *Board) search(history *GameHistory, numberOfMoves int) BoardSearchResults {
	return BoardSearchResults{}
}

func (b *Board) generateMoves() []Move {
	var moves []Move

	return moves
}
