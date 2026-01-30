package engine

import (
	"cmp"
	"fmt"
	"math/bits"
	"slices"
	"strconv"
	"strings"
)

/*
This file contains functionality related to setup, searching and evalution of a board.
*/

// Take a FEN string and turn it into a board, ready for the engine to search over it
// [BUG] FEN logic does not handle castling rights correct, as KQ is valid, but will result in no rights
// [BUG] Should allow KQ, but for now just put in KQ-- and it works
func (position FEN) toBoard(history []FEN) (*Board, error) {
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
	// If not of length 4, then no castling rights
	board.CR = 0
	if len(castling) == 4 {
		// This is more condensed than 4 if statements, but it could accept invalid FENs
		// KQkq is the correct order, but this would accept them out of order, like kQKq
		// This is mostly fine. Maybe in future can just be verbose and check each character one by one
		for s, v := range map[rune]uint8{CHAR_WK: CASTLE_WK, CHAR_WQ: CASTLE_WQ, CHAR_BK: CASTLE_BK, CHAR_BQ: CASTLE_BQ} {
			if strings.Contains(castling, string(s)) {
				board.CR |= v
			}
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

	// Setup the board history if needed
	if history != nil {
		board.History = make([]ZobristHash, 0, STARTING_HISTORY_LENGTH)
	}

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

	// Set all mailboxes to empty, the below loop will populate with the correct pieces
	for i := range NUM_SQUARES {
		b.MailBox[i] = NO_PIECE
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
					b.MailBox[idx] = KING
				case CHAR_BQ:
					b.Pieces[BLACK][QUEEN] |= position
					b.Occupancy[BLACK] |= position
					b.MailBox[idx] = QUEEN
				case CHAR_BR:
					b.Pieces[BLACK][ROOK] |= position
					b.Occupancy[BLACK] |= position
					b.MailBox[idx] = ROOK
				case CHAR_BB:
					b.Pieces[BLACK][BISHOP] |= position
					b.Occupancy[BLACK] |= position
					b.MailBox[idx] = BISHOP
				case CHAR_BN:
					b.Pieces[BLACK][KNIGHT] |= position
					b.Occupancy[BLACK] |= position
					b.MailBox[idx] = KNIGHT
				case CHAR_BP:
					b.Pieces[BLACK][PAWN] |= position
					b.Occupancy[BLACK] |= position
					b.MailBox[idx] = PAWN
				case CHAR_WK:
					b.Pieces[WHITE][KING] |= position
					b.Occupancy[WHITE] |= position
					b.KingSquare[WHITE] = idx
					b.MailBox[idx] = KING
				case CHAR_WQ:
					b.Pieces[WHITE][QUEEN] |= position
					b.Occupancy[WHITE] |= position
					b.MailBox[idx] = QUEEN
				case CHAR_WR:
					b.Pieces[WHITE][ROOK] |= position
					b.Occupancy[WHITE] |= position
					b.MailBox[idx] = ROOK
				case CHAR_WB:
					b.Pieces[WHITE][BISHOP] |= position
					b.Occupancy[WHITE] |= position
					b.MailBox[idx] = BISHOP
				case CHAR_WN:
					b.Pieces[WHITE][KNIGHT] |= position
					b.Occupancy[WHITE] |= position
					b.MailBox[idx] = KNIGHT
				case CHAR_WP:
					b.Pieces[WHITE][PAWN] |= position
					b.Occupancy[WHITE] |= position
					b.MailBox[idx] = PAWN
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

func (b *Board) buildGameHistory(history []FEN) error {
	for _, h := range history {
		b, err := h.toBoard(nil)
		if err != nil {
			return err
		}

		b.History = append(b.History, b.toZobrist())
	}

	return nil
}

type BoardSearchResults struct {
	Nodes     int32
	MoveEvals []MoveEval
}

func (b *Board) search(numberOfMoves int) BoardSearchResults {
	return BoardSearchResults{}
}

// This function generates all legal moves in a position
// DO NOT USE THIS IN THE SEARCH OR ENGINE HOTPATH
// This should only be used for giving the frontend the legal moves in a position
func (b *Board) generateLegalMoves() []Move {
	moves := make([]Move, 0, MAX_NUMBER_OF_MOVES_IN_A_POSITION)
	legalMoves := make([]Move, 0, MAX_NUMBER_OF_MOVES_IN_A_POSITION)
	numberOfMoves := b.generatePseudoLegalMoves(moves)

	for i := range moves[:numberOfMoves] {
		if b.isMoveLegal(moves[i]) {
			legalMoves = append(legalMoves, moves[i])
		}
	}

	return legalMoves
}

// This function generates all pseduo legal moves in a position and fills out a pre-allocated move array
// This is desirable as the engine will check the legality of the move in the search itself
// This could avoid calling isMoveLegal for 30+ moves if we hit an AB-cutoff early, which is a big optimization
func (b *Board) generatePseudoLegalMoves(moves []Move) int {
	// Keep track of where we are in the moves array
	moveIdx := 0

	// Generate pseudo-legal pawn moves
	moveIdx = b.getPawnMoves(moves, moveIdx)

	// Generate pseudo-legal knight moves
	moveIdx = b.getKnightMoves(moves, moveIdx)

	// Generate pseudo-legal king moves
	moveIdx = b.getKingMoves(moves, moveIdx)

	// Generate pseudo-legal bishop moves
	moveIdx = b.getBishopMoves(moves, moveIdx)

	// Generate pseudo-legal rook moves
	moveIdx = b.getRookMoves(moves, moveIdx)

	// Generate pseudo-legal queen moves
	moveIdx = b.getQueenMoves(moves, moveIdx)

	// Generate pseudo-legal castling moves
	moveIdx = b.getCastlingMoves(moves, moveIdx)

	return moveIdx
}

func (b *Board) generatePseudoLegalMovesNegaMax(moves []Move, hasTTEntry bool, ttEntry TTEntry) int {
	moveIdx := b.generatePseudoLegalMoves(moves)

	// Sort the moves to prune more nodes
	slices.SortFunc(moves[:moveIdx], func(ma, mb Move) int {
		return cmp.Compare(mb.orderScore(hasTTEntry, ttEntry.move), ma.orderScore(hasTTEntry, ttEntry.move))
	})

	return moveIdx
}

// This function unmakes a move, in-place, on a board
func (b *Board) unMakeMove(unmove MoveUndo) {

	// Pop a zobrist entry off the history and restore it
	historyLength := len(b.History)
	b.Zobrist = b.History[historyLength-1]
	b.History = b.History[:historyLength-1]

	// Roll back the full move counter
	if b.Turn == WHITE {
		b.FMC--
	}

	// Reset the half move couner
	b.HMC = unmove.hmc

	// Reset EPS
	b.EPS = unmove.eps

	// Reset the color
	b.Turn ^= 1

	// Restore castling rights
	b.CR = unmove.cr

	// Decode move and board
	color := b.Turn
	oppColor := color ^ 1
	start := unmove.start
	target := unmove.target
	captured := unmove.captured
	code := unmove.code
	isPromotion := unmove.isPromotion
	startBitBoard := start.bitBoardPosition()
	targetBitBoard := target.bitBoardPosition()

	// Get piece at target square (the one that moved there from the start square)
	targetPiece := b.getPieceAt(target)

	// Reset bitboard and mailbox
	b.Pieces[color][targetPiece].clear(targetBitBoard)
	b.Occupancy[color].clear(targetBitBoard)
	b.Pieces[color][targetPiece].set(startBitBoard)
	b.Occupancy[color].set(startBitBoard)
	b.MailBox[start] = targetPiece
	b.MailBox[target] = NO_PIECE

	// If piece captured was not NO_PIECE update those bit boards
	if captured != NO_PIECE {

		// If code was en passent, target square is one row off of where peice should go back too
		capturedSq := target
		if code == MOVE_CODE_EN_PASSANT {
			if color == WHITE {
				capturedSq -= 8
			} else {
				capturedSq += 8
			}
		}
		b.Pieces[oppColor][captured].set(capturedSq.bitBoardPosition())
		b.Occupancy[oppColor].set(capturedSq.bitBoardPosition())
		b.MailBox[capturedSq] = captured
	}

	// Handle promotion by downgrading back to a pawn
	if isPromotion {
		b.Pieces[color][targetPiece].clear(startBitBoard)
		b.Pieces[color][PAWN].set(startBitBoard)
		b.MailBox[start] = PAWN
	}

	// Handle castling
	// Castling rights are already restored above, this is for moving the rook back into the corner
	if code == MOVE_CODE_CASTLE {
		// White Kingside
		if target == G1 {
			h1BB := H1.bitBoardPosition()
			f1BB := F1.bitBoardPosition()
			b.Pieces[WHITE][ROOK].clear(f1BB)
			b.Occupancy[WHITE].clear(f1BB)
			b.Pieces[WHITE][ROOK].set(h1BB)
			b.Occupancy[WHITE].set(h1BB)
			b.MailBox[F1] = NO_PIECE
			b.MailBox[H1] = ROOK
		}

		// White Queenside
		if target == C1 {
			a1BB := A1.bitBoardPosition()
			d1BB := D1.bitBoardPosition()
			b.Pieces[WHITE][ROOK].clear(d1BB)
			b.Occupancy[WHITE].clear(d1BB)
			b.Pieces[WHITE][ROOK].set(a1BB)
			b.Occupancy[WHITE].set(a1BB)
			b.MailBox[D1] = NO_PIECE
			b.MailBox[A1] = ROOK
		}

		// Black Kingside
		if target == G8 {
			h8BB := H8.bitBoardPosition()
			f8BB := F8.bitBoardPosition()
			b.Pieces[BLACK][ROOK].clear(f8BB)
			b.Occupancy[BLACK].clear(f8BB)
			b.Pieces[BLACK][ROOK].set(h8BB)
			b.Occupancy[BLACK].set(h8BB)
			b.MailBox[F8] = NO_PIECE
			b.MailBox[H8] = ROOK
		}

		// Black Queenside
		if target == C8 {
			a8BB := A8.bitBoardPosition()
			d8BB := D8.bitBoardPosition()
			b.Pieces[BLACK][ROOK].clear(d8BB)
			b.Occupancy[BLACK].clear(d8BB)
			b.Pieces[BLACK][ROOK].set(a8BB)
			b.Occupancy[BLACK].set(a8BB)
			b.MailBox[D8] = NO_PIECE
			b.MailBox[A8] = ROOK
		}
	}

	// Update king square
	if targetPiece == KING {
		b.KingSquare[color] = start
	}

	// Update either color occupancy
	b.Occupancy[EITHER_COLOR] = (b.Occupancy[WHITE] | b.Occupancy[BLACK])
}

// This function makes a move, in-place, on a board, and returns if that move was legal or not
func (b *Board) makeMove(move Move) (MoveUndo, bool) {

	// Create an unmake entry somewhere
	unmake := MoveUndo{
		hmc:         b.HMC,
		cr:          b.CR,
		eps:         b.EPS,
		captured:    NO_PIECE,
		isPromotion: false,
	}

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

	// Put decoded move information into the unmake struct
	unmake.start = start
	unmake.target = target
	unmake.code = code

	// Hash out enpassent square
	if b.EPS != NO_SQUARE {
		b.Zobrist ^= ENPASSENT_ZOBRIST[b.EPS%8]
	}

	// Reset enpassent (will get reset later if needed)
	b.EPS = NO_SQUARE

	// Get pieces of start and target squares
	startPiece := b.getPieceAt(start)

	// Handle captures
	if code == MOVE_CODE_CAPTURE {
		targetPiece := b.getPieceAt(target)

		// Put the captured piece into the unmake struct
		unmake.captured = targetPiece

		// Clear enemy piece and reset half move clock
		b.Pieces[oppColor][targetPiece].clear(targetBitBoard)
		b.Occupancy[oppColor].clear(targetBitBoard)
		b.HMC = 0

		// Hash out enempy piece
		b.Zobrist ^= PIECE_ZOBRIST[oppColor][targetPiece][target]
	}

	// Handle the moving piece bitboards
	b.Pieces[color][startPiece].clear(startBitBoard)
	b.Occupancy[color].clear(startBitBoard)
	b.Pieces[color][startPiece].set(targetBitBoard)
	b.Occupancy[color].set(targetBitBoard)

	// Hash out old piece and in new piece
	b.Zobrist ^= PIECE_ZOBRIST[color][startPiece][start]
	b.Zobrist ^= PIECE_ZOBRIST[color][startPiece][target]

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

			// Hash out captured enemy pawn
			b.Zobrist ^= PIECE_ZOBRIST[oppColor][PAWN][eps]

			// Add to unmake struct
			unmake.captured = PAWN
		}

		// Handle pawn promoting
		if promotion != NO_PIECE {
			b.Pieces[color][PAWN].clear(targetBitBoard)
			b.Pieces[color][promotion].set(targetBitBoard)
			b.MailBox[target] = promotion

			// Hash out pawn and in promoted piece
			b.Zobrist ^= PIECE_ZOBRIST[color][PAWN][target]
			b.Zobrist ^= PIECE_ZOBRIST[color][promotion][target]

			// Store it for unmake later
			unmake.isPromotion = true
		}
	}

	// Hash in enpassent square (if not none)
	if b.EPS != NO_SQUARE {
		b.Zobrist ^= ENPASSENT_ZOBRIST[b.EPS%8]
	}

	// Hash out castling rights
	b.Zobrist ^= CASTLING_ZOBRIST[b.CR]

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

			// Update Zobrist hash for the rook
			b.Zobrist ^= PIECE_ZOBRIST[WHITE][ROOK][H1]
			b.Zobrist ^= PIECE_ZOBRIST[WHITE][ROOK][F1]

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

			// Update Zobrist hash for the rook
			b.Zobrist ^= PIECE_ZOBRIST[WHITE][ROOK][A1]
			b.Zobrist ^= PIECE_ZOBRIST[WHITE][ROOK][D1]

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

			// Update Zobrist hash for the rook
			b.Zobrist ^= PIECE_ZOBRIST[BLACK][ROOK][H8]
			b.Zobrist ^= PIECE_ZOBRIST[BLACK][ROOK][F8]

			// Clear black castling rights
			b.CR &= ^uint8(CASTLE_BK)
			b.CR &= ^uint8(CASTLE_BQ)
		}

		// Black Queenside
		if target == C8 {
			a8BB := A8.bitBoardPosition()
			d8BB := D8.bitBoardPosition()
			b.Pieces[BLACK][ROOK].clear(a8BB)
			b.Occupancy[BLACK].clear(a8BB)
			b.Pieces[BLACK][ROOK].set(d8BB)
			b.Occupancy[BLACK].set(d8BB)
			b.MailBox[A8] = NO_PIECE
			b.MailBox[D8] = ROOK

			// Update Zobrist hash for the rook
			b.Zobrist ^= PIECE_ZOBRIST[BLACK][ROOK][A8]
			b.Zobrist ^= PIECE_ZOBRIST[BLACK][ROOK][D8]

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

	// Hash in castling rights
	b.Zobrist ^= CASTLING_ZOBRIST[b.CR]

	// Update board turn and update either color occupancy
	b.Turn ^= 1
	b.Occupancy[EITHER_COLOR] = b.Occupancy[WHITE] | b.Occupancy[BLACK]

	// Hash in turn
	b.Zobrist ^= BLACK_TO_MOVE_ZOBRIST

	// Verify the board start is legal
	// Make sure the king is not attacked
	if b.isSquareAttacked(b.KingSquare[color], b.Turn) {
		return unmake, false
	}

	return unmake, true
}

// This function simply checks if a move was legal, utilizing make and unmake moves
func (b *Board) isMoveLegal(move Move) bool {
	unmove, isLegal := b.makeMove(move)
	b.unMakeMove(unmove)
	return isLegal
}

// This function returns the piece at a specific square
func (b *Board) getPieceAt(sq Square) Piece {
	return b.MailBox[sq]
}

// This function prints the board to stdout, useful for debugging or CLI
func (b *Board) print() {
	for i := 7; i >= 0; i-- {
		for j := range 8 {
			sq := Square(i*8 + j)
			piece := b.getPieceAt(sq)
			color := WHITE
			if b.Occupancy[WHITE]&sq.bitBoardPosition() == 0 {
				color = BLACK
			}
			fmt.Print(piece.toString(color) + " ")
		}
		fmt.Println()
	}
}

// Useful for debugging move and unmove
func (b *Board) desync(where string) {
	for sq := range NUM_SQUARES {
		sqBB := Square(sq).bitBoardPosition()

		// Get Mailbox State
		mbPiece := b.getPieceAt(Square(sq))
		mbColor := WHITE
		if b.Occupancy[BLACK]&sqBB != 0 {
			mbColor = BLACK
		}

		// Get Bitboard State
		bbPiece := NO_PIECE
		bbColor := WHITE

		// Check White Pieces
		if b.Pieces[WHITE][KING]&sqBB != 0 {
			bbPiece = KING
			bbColor = WHITE
		} else if b.Pieces[WHITE][QUEEN]&sqBB != 0 {
			bbPiece = QUEEN
			bbColor = WHITE
		} else if b.Pieces[WHITE][ROOK]&sqBB != 0 {
			bbPiece = ROOK
			bbColor = WHITE
		} else if b.Pieces[WHITE][BISHOP]&sqBB != 0 {
			bbPiece = BISHOP
			bbColor = WHITE
		} else if b.Pieces[WHITE][KNIGHT]&sqBB != 0 {
			bbPiece = KNIGHT
			bbColor = WHITE
		} else if b.Pieces[WHITE][PAWN]&sqBB != 0 {
			bbPiece = PAWN
			bbColor = WHITE
		} else if b.Pieces[BLACK][KING]&sqBB != 0 {
			bbPiece = KING
			bbColor = BLACK
		} else if b.Pieces[BLACK][QUEEN]&sqBB != 0 {
			bbPiece = QUEEN
			bbColor = BLACK
		} else if b.Pieces[BLACK][ROOK]&sqBB != 0 {
			bbPiece = ROOK
			bbColor = BLACK
		} else if b.Pieces[BLACK][BISHOP]&sqBB != 0 {
			bbPiece = BISHOP
			bbColor = BLACK
		} else if b.Pieces[BLACK][KNIGHT]&sqBB != 0 {
			bbPiece = KNIGHT
			bbColor = BLACK
		} else if b.Pieces[BLACK][PAWN]&sqBB != 0 {
			bbPiece = PAWN
			bbColor = BLACK
		}

		// Compare Piece Types
		if mbPiece != bbPiece {
			b.print()
			fmt.Printf("DESYNC AT %v (%s) -> Mailbox: %v, Bitboard: %v\n", Square(sq).toString(), where, mbPiece.toString(mbColor), bbPiece.toString(bbColor))
			panic("Board desync: Piece Mismatch")
		}

		// Compare Colors (Only if a piece actually exists)
		if mbPiece != NO_PIECE {

			if mbColor != bbColor {
				b.print()
				fmt.Printf("DESYNC AT %d (%s) -> Mailbox Color: %d, Bitboard Color: %d\n", sq, where, mbColor, bbColor)
				panic("Board desync: Color Mismatch")
			}
		}
	}
}

// Function used to get the phase score of the position
// This should be called once at the start of the evaluation of a board
// This returns the phase offset to be used against the PST tables
// Function used by the board to get the pst value
func (b *Board) getPhaseScore() int {
	// Weights for calculating game phase (Non-Pawn Material is standard)
	const (
		PawnPhase   = 0
		KnightPhase = 1
		BishopPhase = 1
		RookPhase   = 2
		QueenPhase  = 4
	)

	// Calculate the maximum possible phase (Starting Position)
	// 16 Pawns, 4 Knights, 4 Bishops, 4 Rooks, 2 Queens
	const TotalPhase = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

	// Start with TotalPhase and subtract the pieces currently on the board.
	// If board is full -> phase = TotalPhase - TotalPhase = 0 (Opening)
	// If board is empty -> phase = TotalPhase - 0 = TotalPhase (Endgame)
	phase := TotalPhase

	// Count White Pieces
	wp := bits.OnesCount64(uint64(b.Pieces[WHITE][PAWN]))
	wn := bits.OnesCount64(uint64(b.Pieces[WHITE][KNIGHT]))
	wb := bits.OnesCount64(uint64(b.Pieces[WHITE][BISHOP]))
	wr := bits.OnesCount64(uint64(b.Pieces[WHITE][ROOK]))
	wq := bits.OnesCount64(uint64(b.Pieces[WHITE][QUEEN]))

	// Count Black Pieces
	bp := bits.OnesCount64(uint64(b.Pieces[BLACK][PAWN]))
	bn := bits.OnesCount64(uint64(b.Pieces[BLACK][KNIGHT]))
	bb := bits.OnesCount64(uint64(b.Pieces[BLACK][BISHOP]))
	br := bits.OnesCount64(uint64(b.Pieces[BLACK][ROOK]))
	bq := bits.OnesCount64(uint64(b.Pieces[BLACK][QUEEN]))

	// Subtract material currently on board
	phase -= wp * PawnPhase
	phase -= bp * PawnPhase
	phase -= wn * KnightPhase
	phase -= bn * KnightPhase
	phase -= wb * BishopPhase
	phase -= bb * BishopPhase
	phase -= wr * RookPhase
	phase -= br * RookPhase
	phase -= wq * QueenPhase
	phase -= bq * QueenPhase

	// Normalize to range [0, 256]
	// 0   = Opening
	// 256 = Endgame
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase

	return phase
}
