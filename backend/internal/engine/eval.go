package engine

import "math/bits"

/*
This file holds all the functionality related to the evaluation of a static board
*/

// Defining hard material values for each piece, just raw not sophisticated way of looking at material
const (
	PAWN_VALUE   Eval = 100
	KNIGHT_VALUE Eval = 300
	BISHOP_VALUE Eval = 300
	ROOK_VALUE   Eval = 500
	QUEEN_VALUE  Eval = 900
	KING_VALUE   Eval = 20000
)

var PIECE_VALUES [NUM_PIECES]Eval

// File masks used for fast evaluation
var FileMask [8]BitBoard

// Passed pawn masks for fast evaluation
var WhitePassedMask [64]BitBoard
var BlackPassedMask [64]BitBoard

// King safety pask for fast evaluation
var KingSafetyMask [64]BitBoard

// Init evaluation, called at engine startup
func initEval() {
	PIECE_VALUES[PAWN] = PAWN_VALUE
	PIECE_VALUES[KNIGHT] = KNIGHT_VALUE
	PIECE_VALUES[BISHOP] = BISHOP_VALUE
	PIECE_VALUES[ROOK] = ROOK_VALUE
	PIECE_VALUES[QUEEN] = QUEEN_VALUE
	PIECE_VALUES[KING] = KING_VALUE

	// Init file mask
	for file := range 8 {
		for rank := range 8 {
			sq := rank*8 + file
			FileMask[file] |= 1 << sq
		}
	}

	// Init passed pawn mask
	for sq := range NUM_SQUARES {
		f := sq % 8
		r := sq / 8

		for target := range NUM_SQUARES {
			tf := target % 8
			tr := target / 8
			if (tf >= f-1 && tf <= f+1) && (tr > r) {
				WhitePassedMask[sq] |= (BitBoard(1) << target)
			}
			if (tf >= f-1 && tf <= f+1) && (tr < r) {
				BlackPassedMask[sq] |= (BitBoard(1) << target)
			}
		}
	}

	// Init king safety mask
	for sq := range NUM_SQUARES {
		f := sq % 8
		r := sq / 8

		for df := -2; df <= 2; df++ {
			for dr := -2; dr <= 2; dr++ {
				// Skip the square the king is actually standing on
				if df == 0 && dr == 0 {
					continue
				}

				tf, tr := f+df, r+dr
				if tf >= 0 && tf < 8 && tr >= 0 && tr < 8 {
					KingSafetyMask[sq] |= (BitBoard(1) << (tr*8 + tf))
				}
			}
		}
	}
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
	return interpolatePhase(phaseSocre, openingEval, endgameEval)
}

// This function will take the phase score and two evals and return the interpolated score
func interpolatePhase(phaseScore int, opneing, endgame Eval) Eval {
	return Eval(((int(opneing) * (256 - phaseScore)) + (int(endgame) * phaseScore)) / 256)
}

// Function to get the evalution based on the pawn structure of the board
func (b *Board) pawnStructureEval(phaseScore int) Eval {
	eval := Eval(0)

	// Get pawn bitboards
	whitePawns := b.Pieces[WHITE][PAWN]
	blackPawns := b.Pieces[BLACK][PAWN]

	// Doubled pawns
	// Each doubled pawn will be worth a 15 centipawn penalty to 36 point penalty
	// The penalty increases as the game moves towards the engame
	doubledPenalty := interpolatePhase(phaseScore, 15, 36)
	doubledWhite := doubledPawns(whitePawns)
	doubledBlack := doubledPawns(blackPawns)
	eval -= Eval(bits.OnesCount64(uint64(doubledWhite))) * doubledPenalty
	eval += Eval(bits.OnesCount64(uint64(doubledBlack))) * doubledPenalty

	// Isolated pawns
	// Each isolated pawn will be worth a 19 centipawn penalty to 31 point penalty
	// The penalty increases as the game moves towards the engame
	isolatedPenalty := interpolatePhase(phaseScore, 19, 31)
	isolatedWhite := isolatedPawns(whitePawns)
	isolatedBlack := isolatedPawns(blackPawns)
	eval -= Eval(bits.OnesCount64(uint64(isolatedWhite))) * isolatedPenalty
	eval += Eval(bits.OnesCount64(uint64(isolatedBlack))) * isolatedPenalty

	// Passed pawns
	// Each passed pawn will be worth a 15 centipawn bonus to 38 point bonus
	// Potentially in the future make this row dependent scores
	passedBonus := interpolatePhase(phaseScore, 15, 38)
	passedWhite := passedPawnsWhite(whitePawns, blackPawns)
	passedBlack := passedPawnsBlack(blackPawns, whitePawns)
	eval += Eval(passedWhite) * passedBonus
	eval -= Eval(passedBlack) * passedBonus

	return eval
}

// Get white passed pawns
func passedPawnsWhite(whitePawns, blackPawns BitBoard) int {
	count := 0
	for whitePawns > 0 {
		sq := whitePawns.popSquare()
		if WhitePassedMask[sq]&blackPawns == 0 {
			count++
		}
	}

	return count
}

// Get black passed pawns
func passedPawnsBlack(blackPawns, whitePawns BitBoard) int {
	count := 0
	for blackPawns > 0 {
		sq := blackPawns.popSquare()
		if BlackPassedMask[sq]&whitePawns == 0 {
			count++
		}
	}

	return count
}

// Function to get doubled pawns
func doubledPawns(pawns BitBoard) int {
	doubled := 0

	for file := range 8 {
		onFile := pawns & FileMask[file]
		if bits.OnesCount64(uint64(onFile)) >= 2 {
			doubled++
		}
	}
	return doubled
}

// Function to get isolated pawns
func isolatedPawns(pawns BitBoard) BitBoard {
	isolated := BitBoard(0)

	for file := range 8 {
		onFile := pawns & FileMask[file]

		// Adjacent files
		var neighbors BitBoard
		if file > 0 {
			neighbors |= pawns & FileMask[file-1]
		}
		if file < 7 {
			neighbors |= pawns & FileMask[file+1]
		}

		// Isolated pawns on this file
		isolated |= onFile &^ neighbors
	}

	return isolated
}

// Function to evaluate king safety
func (b *Board) kingSafetyEval(phaseScore int) Eval {
	eval := Eval(0)

	// Get phase scores
	friendlyPawnScore := interpolatePhase(phaseScore, 17, 0)
	enemyPawnScore := interpolatePhase(phaseScore, 13, 0)
	friendlyPieceScore := interpolatePhase(phaseScore, 7, 0)
	enemyPieceScore := interpolatePhase(phaseScore, 13, 0)

	// Do white king safety
	// Get the saftey mask, and check for white pawns/pieces and black pawns/pieces
	// Evaluate strong for number of pawns in front of king, and harshly for enemy pawns/pieces next to king
	safetyMask := KingSafetyMask[b.KingSquare[WHITE]]
	friendlyPawns := b.Pieces[WHITE][PAWN] & safetyMask
	enemyPawns := b.Pieces[BLACK][PAWN] & safetyMask
	friendlyPieces := (b.Occupancy[WHITE] &^ friendlyPawns) & safetyMask
	enemyPieces := (b.Occupancy[BLACK] &^ enemyPawns) & safetyMask
	eval += Eval(bits.OnesCount64(uint64(friendlyPawns))) * friendlyPawnScore
	eval += Eval(bits.OnesCount64(uint64(friendlyPieces))) * friendlyPieceScore
	eval -= Eval(bits.OnesCount64(uint64(enemyPawns))) * enemyPawnScore
	eval -= Eval(bits.OnesCount64(uint64(enemyPieces))) * enemyPieceScore

	// Do black king safety
	// Get the saftey mask, and check for black pawns/pieces and white pawns/pieces
	// Evaluate strong for number of pawns in front of king, and harshly for enemy pawns/pieces next to king
	safetyMask = KingSafetyMask[b.KingSquare[BLACK]]
	friendlyPawns = b.Pieces[BLACK][PAWN] & safetyMask
	enemyPawns = b.Pieces[WHITE][PAWN] & safetyMask
	friendlyPieces = (b.Occupancy[BLACK] &^ friendlyPawns) & safetyMask
	enemyPieces = (b.Occupancy[WHITE] &^ enemyPawns) & safetyMask
	eval -= Eval(bits.OnesCount64(uint64(friendlyPawns))) * friendlyPawnScore
	eval -= Eval(bits.OnesCount64(uint64(friendlyPieces))) * friendlyPieceScore
	eval += Eval(bits.OnesCount64(uint64(enemyPawns))) * enemyPawnScore
	eval += Eval(bits.OnesCount64(uint64(enemyPieces))) * enemyPieceScore

	return eval
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

	// Do pawn structure eval
	eval += b.pawnStructureEval(phaseSocre)

	// Do king safety
	eval += b.kingSafetyEval(phaseSocre)

	return eval
}
