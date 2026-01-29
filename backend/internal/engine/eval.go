package engine

/*
This file holds all the functionality related to the evaluation of a static board
*/

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
	finalEval := Eval(((int(openingEval) * (256 - phaseSocre)) + (int(endgameEval) * phaseSocre)) / 256)

	return finalEval
}

// Main evaluation function, to be called by the searching algorithm
func (b *Board) eval() Eval {
	// Get the current phase of the board
	phaseSocre := b.getPhaseScore()

	// Simple pst evaluation
	return b.pstEval(phaseSocre)
}
