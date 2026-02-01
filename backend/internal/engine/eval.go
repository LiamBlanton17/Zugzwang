package engine

import "math/bits"

/*
This file holds all the functionality related to the evaluation of a static board
*/

// Defining hard material values for each piece, just raw not sophisticated way of looking at material
const (
	PAWN_VALUE   Eval = 100
	KNIGHT_VALUE Eval = 300
	BISHOP_VALUE Eval = 310
	ROOK_VALUE   Eval = 500
	QUEEN_VALUE  Eval = 900
	KING_VALUE   Eval = 20000
)

var PIECE_VALUES [NUM_PIECES]Eval

// Init evaluation, called at engine startup
func initEval() {
	PIECE_VALUES[PAWN] = PAWN_VALUE
	PIECE_VALUES[KNIGHT] = KNIGHT_VALUE
	PIECE_VALUES[BISHOP] = BISHOP_VALUE
	PIECE_VALUES[ROOK] = ROOK_VALUE
	PIECE_VALUES[QUEEN] = QUEEN_VALUE
	PIECE_VALUES[KING] = KING_VALUE
}

// Piece square table evaluation of the position
func (b *Board) pstEval(phaseSocre int) Eval {
	openingEval := Eval(0)
	endgameEval := Eval(0)

	for p := PAWN; p <= KING; p++ {
		// Add White Pieces
		whiteBitboard := b.Pieces[WHITE][p]
		for whiteBitboard != 0 {
			sq := whiteBitboard.popSquare()
			openingEval += PST[OPENING][WHITE][p][sq]
			endgameEval += PST[ENDGAME][WHITE][p][sq]
		}

		// Subtract Black Pieces
		blackBitboard := b.Pieces[BLACK][p]
		for blackBitboard != 0 {
			sq := blackBitboard.popSquare()
			openingEval -= PST[OPENING][BLACK][p][sq]
			endgameEval -= PST[ENDGAME][BLACK][p][sq]
		}
	}

	// Interperlate with the phase score to get the final eval
	return Eval(((int(openingEval) * (256 - phaseSocre)) + (int(endgameEval) * phaseSocre)) / 256)
}

// Main evaluation function, to be called by the searching algorithm
func (b *Board) eval() Eval {
	// Get the current phase of the board
	phaseSocre := b.getPhaseScore()

	// Simple pst evaluation
	eval := b.pstEval(phaseSocre)

	// Simple tempo evaluation
	if b.Turn == WHITE {
		eval += 10
	} else {
		eval -= 10
	}

	// Simple bishop pair evaluation
	if bits.OnesCount64(uint64(b.Pieces[WHITE][BISHOP])) >= 2 {
		eval += 30
	}
	if bits.OnesCount64(uint64(b.Pieces[BLACK][BISHOP])) >= 2 {
		eval -= 30
	}

	return eval
}
