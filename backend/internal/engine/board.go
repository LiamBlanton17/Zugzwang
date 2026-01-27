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
					b.KingSquare[BLACK] = idx
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
					b.KingSquare[WHITE] = idx
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
	return b.Occupancy[b.Turn^1]
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

// This function generates all legal moves in a position
// DO NOT USE THIS IN THE SEARCH OR ENGINE HOTPATH
// This should only be used for giving the frontend the legal moves in a position
func (b *Board) generateLegalMoves() []Move {
	b.generatePseudoLegalMoves()
	moves := b.Moves
	legalMoves := make([]Move, 0, len(b.Moves))

	for i := range b.MoveIdx {
		if b.isMoveLegal(moves[i]) {
			legalMoves = append(legalMoves, b.Moves[i])
		}
	}

	return legalMoves
}

// This function generates all pseduo legal moves in a position and fills out a pre-allocated move array
// This is desirable as the engine will check the legality of the move in the search itself
// This could avoid calling isMoveLegal for 30+ moves if we hit an AB-cutoff early, which is a big optimization
func (b *Board) generatePseudoLegalMoves() {
	// Reset board move index
	b.MoveIdx = 0

	// Generate pseudo-legal pawn moves
	b.getPawnMoves()

	// Generate pseudo-legal knight moves
	b.getKnightMoves()

	// Generate pseudo-legal king moves
	b.getKingMoves()

	// Generate pseudo-legal bishop moves
	b.getBishopMoves()

	// Generate pseudo-legal rook moves
	b.getRookMoves()

	// Generate pseudo-legal queen moves
	b.getQueenMoves()

	// Generate pseudo-legal castling moves
	b.getCastlingMoves()
}

// This function makes a move, in-place, on a board, and returns if that move was legal or not
func (b *Board) makeMove(move Move) bool {

	// Create an unmake entry somewhere

	// Add this boards Zobrist hash to the history and update clocks
	b.History = append(b.History, b.Zobrist)
	b.HMC++
	if b.Turn == BLACK {
		b.FMC++
	}

	// Decode the move and board start
	start := move.start
	target := move.target
	startBitBoard := start.bitBoardPosition()
	targetBitBoard := target.bitBoardPosition()
	promotion := move.promotion
	code := move.code
	color := b.Turn
	oppColor := color ^ 1
	eps := b.EPS

	// Reset enpassent (will get reset later if needed)
	b.EPS = NO_SQUARE

	// Get pieces of start and target squares
	startPiece := b.getPieceAt(start)

	// Handle captures
	if code == MOVE_CODE_CAPTURE {
		targetPiece := b.getPieceAt(target)

		// Clear enemy piece and reset half move clock
		b.Pieces[oppColor][targetPiece].clear(targetBitBoard)
		b.Occupancy[oppColor].clear(targetBitBoard)
		b.HMC = 0
	}

	// Handle the moving piece bitboards
	b.Pieces[color][startPiece].clear(startBitBoard)
	b.Occupancy[color].clear(startBitBoard)
	b.Pieces[color][startPiece].set(targetBitBoard)
	b.Occupancy[color].set(targetBitBoard)

	// Handle the piece mailbox
	b.MailBox[start] = NO_PIECE
	b.MailBox[target] = startPiece

	// Handle pawns and en passent
	if startPiece == PAWN {
		b.HMC = 0

		// Handle pawn being double pushed
		// Update the en passent square of the board
		if code == MOVE_CODE_DOUBLE_PAWN_PUSH {
			b.EPS = target
			if color == WHITE {
				b.EPS -= 8
			} else {
				b.EPS += 8
			}
		}

		// Handle pawn capturing en passant
		if code == MOVE_CODE_EN_PASSANT {
			if color == WHITE {
				eps -= 8
			} else {
				eps += 8
			}
			epsBitBoard := eps.bitBoardPosition()
			b.Pieces[oppColor][PAWN].clear(epsBitBoard)
			b.Occupancy[oppColor].clear(epsBitBoard)
			b.MailBox[eps] = NO_PIECE
		}

		// Handle pawn promoting
		if code == MOVE_CODE_PROMOTION {
			b.Pieces[color][PAWN].clear(targetBitBoard)
			b.Pieces[color][promotion].set(targetBitBoard)
			b.MailBox[target] = promotion
		}
	}

	// Handle castling (moving rook and updating castling rights)
	if code == MOVE_CODE_CASTLE {
		// White Kingside
		if target == G1 {
			h1BB := H1.bitBoardPosition()
			f1BB := F1.bitBoardPosition()
			b.Pieces[WHITE][ROOK].clear(h1BB)
			b.Occupancy[WHITE].clear(h1BB)
			b.Pieces[WHITE][ROOK].set(f1BB)
			b.Occupancy[WHITE].set(f1BB)
			b.MailBox[H1] = NO_PIECE
			b.MailBox[F1] = ROOK

			// Clear white castling rights
			b.CR &= ^uint8(CASTLE_WK)
			b.CR &= ^uint8(CASTLE_WQ)
		}

		// White Queenside
		if target == C1 {
			a1BB := A1.bitBoardPosition()
			d1BB := D1.bitBoardPosition()
			b.Pieces[WHITE][ROOK].clear(a1BB)
			b.Occupancy[WHITE].clear(a1BB)
			b.Pieces[WHITE][ROOK].set(d1BB)
			b.Occupancy[WHITE].set(d1BB)
			b.MailBox[A1] = NO_PIECE
			b.MailBox[D1] = ROOK

			// Clear white castling rights
			b.CR &= ^uint8(CASTLE_WK)
			b.CR &= ^uint8(CASTLE_WQ)
		}

		// Black Kingside
		if target == G8 {
			h8BB := H8.bitBoardPosition()
			f8BB := F8.bitBoardPosition()
			b.Pieces[BLACK][ROOK].clear(h8BB)
			b.Occupancy[BLACK].clear(h8BB)
			b.Pieces[BLACK][ROOK].set(f8BB)
			b.Occupancy[BLACK].set(f8BB)
			b.MailBox[H8] = NO_PIECE
			b.MailBox[F8] = ROOK

			// Clear black castling rights
			b.CR &= ^uint8(CASTLE_BK)
			b.CR &= ^uint8(CASTLE_BQ)
		}

		// Black Kingside
		if target == C8 {
			a8BB := A8.bitBoardPosition()
			d8BB := D8.bitBoardPosition()
			b.Pieces[BLACK][ROOK].clear(a8BB)
			b.Occupancy[BLACK].clear(a8BB)
			b.Pieces[BLACK][ROOK].set(d8BB)
			b.Occupancy[BLACK].set(d8BB)
			b.MailBox[A8] = NO_PIECE
			b.MailBox[D8] = ROOK

			// Clear black castling rights
			b.CR &= ^uint8(CASTLE_BK)
			b.CR &= ^uint8(CASTLE_BQ)
		}
	}

	// Update king square
	if startPiece == KING {
		b.KingSquare[color] = target

		// Clear castling rights
		if color == WHITE {
			// Clear white castling rights
			b.CR &= ^uint8(CASTLE_WK)
			b.CR &= ^uint8(CASTLE_WQ)
		} else {
			// Clear black castling rights
			b.CR &= ^uint8(CASTLE_BK)
			b.CR &= ^uint8(CASTLE_BQ)
		}
	}

	// If still has castling rights, check corners for rooks
	// If corner doesn't have rooks (moved or captures) unset castling rights
	if b.CR&CASTLE_WK != 0 || b.CR&CASTLE_WQ != 0 {
		// If rook not on H1, remove white queenside castling rights
		if b.Pieces[WHITE][ROOK]&H1.bitBoardPosition() == 0 {
			b.CR &= ^uint8(CASTLE_WK)
		}
		// If rook not on A1, remove white queenside castling rights
		if b.Pieces[WHITE][ROOK]&A1.bitBoardPosition() == 0 {
			b.CR &= ^uint8(CASTLE_WQ)
		}
	}
	if b.CR&CASTLE_BK != 0 || b.CR&CASTLE_BQ != 0 {
		// If rook not on H1, remove white queenside castling rights
		if b.Pieces[BLACK][ROOK]&H8.bitBoardPosition() == 0 {
			b.CR &= ^uint8(CASTLE_BK)
		}
		// If rook not on A8, remove black queenside castling rights
		if b.Pieces[BLACK][ROOK]&A8.bitBoardPosition() == 0 {
			b.CR &= ^uint8(CASTLE_BQ)
		}
	}

	// Update board turn and update either color occupancy
	b.Turn ^= 1
	b.Occupancy[EITHER_COLOR] = (b.Occupancy[WHITE] | b.Occupancy[BLACK])

	// Verify the board start is legal
	// Make sure the king is not attacked
	if b.isSquareAttacked(b.KingSquare[color], b.Turn) {
		return false
	}

	return true
}

// This function unmakes a move, in-place, on a board
func (b *Board) unMakeMove(move Move) {

}

// This function simply checks if a move was legal, utilizing make and unmake moves
func (b *Board) isMoveLegal(move Move) bool {
	isLegal := b.makeMove(move)
	b.unMakeMove(move)
	return isLegal
}

// This function returns the piece at a specific square
func (b *Board) getPieceAt(sq Square) Piece {
	return b.MailBox[sq]
}
