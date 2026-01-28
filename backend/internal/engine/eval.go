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

// "Dumb" material eval of the poisiton
// This eval does not consider what square the piece is on but just its flat material worth
func (b *Board) dumbMaterialEval() Eval {
	eval := Eval(0)

	for s := range NUM_SQUARES {
		sq := Square(s)
		piece := b.getPieceAt(sq)
		if piece == NO_PIECE {
			continue
		}

		if sq.bitBoardPosition()&b.Occupancy[WHITE] == 0 {
			eval += DUMB_CENTIPAWN[piece]
		} else {
			eval -= DUMB_CENTIPAWN[piece]
		}
	}

	return eval
}

// Main evaluation function, to be called by the searching algorithm
func (b *Board) eval() Eval {
	return b.dumbMaterialEval()
}
