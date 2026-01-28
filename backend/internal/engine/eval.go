package engine

/*
This file holds all the functionality related to the evaluation of a static board
*/

// "Dumb" piece values
var DUMB_CENTIPAWN = [NUM_PIECES]Eval{
	PAWN:   100,
	KNIGHT: 300,
	BISHOP: 320,
	ROOK:   500,
	QUEEN:  900,
	KING:   50000,
}

// Piece square table evaluation of the position
func (b *Board) pstEval() Eval {
	eval := Eval(0)

	for p := PAWN; p <= KING; p++ {
		pieceValue := DUMB_CENTIPAWN[p]

		// 1. Add White Pieces
		whiteBitboard := b.Pieces[WHITE][p]
		for whiteBitboard != 0 {
			sq := whiteBitboard.popSquare()
			eval += PST[OPENING][WHITE][p][sq] + pieceValue
		}

		// 2. Subtract Black Pieces
		blackBitboard := b.Pieces[BLACK][p]
		for blackBitboard != 0 {
			sq := blackBitboard.popSquare()
			eval -= (PST[OPENING][BLACK][p][sq] + pieceValue)
		}
	}

	if b.Turn == BLACK {
		return -eval
	}

	return eval
}

// Main evaluation function, to be called by the searching algorithm
func (b *Board) eval() Eval {
	return b.pstEval()
}
