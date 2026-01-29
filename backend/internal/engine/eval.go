package engine

/*
This file holds all the functionality related to the evaluation of a static board
*/

// Piece square table evaluation of the position
func (b *Board) pstEval() Eval {
	eval := Eval(0)

	for p := PAWN; p <= KING; p++ {
		// Add White Pieces
		whiteBitboard := b.Pieces[WHITE][p]
		for whiteBitboard != 0 {
			sq := whiteBitboard.popSquare()
			eval += PST[OPENING][WHITE][p][sq]
		}

		// Subtract Black Pieces
		blackBitboard := b.Pieces[BLACK][p]
		for blackBitboard != 0 {
			sq := blackBitboard.popSquare()
			eval -= PST[OPENING][BLACK][p][sq]
		}
	}

	return eval
}

// Main evaluation function, to be called by the searching algorithm
func (b *Board) eval() Eval {
	return b.pstEval()
}
