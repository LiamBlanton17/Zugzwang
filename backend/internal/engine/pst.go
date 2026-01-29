package engine

/*
This file holds the global Piece-Square Tables (psts) for the evaluation function
*/

// Base tables

// PAWNS
var tablePawn = [NUM_SQUARES]Eval{
	000, 000, 000, 000, 000, 000, 000, 000,
	100, 100, 100, 90, 90, 100, 100, 100,
	105, 105, 90, 100, 100, 90, 105, 105,
	105, 105, 100, 125, 125, 100, 105, 105,
	110, 110, 105, 135, 135, 105, 110, 110,
	115, 120, 125, 145, 145, 125, 120, 115,
	150, 150, 150, 150, 150, 150, 150, 150,
	000, 000, 000, 000, 000, 000, 000, 000,
}

// KNIGHTS
var tableKnight = [NUM_SQUARES]Eval{
	245, 255, 260, 260, 260, 260, 255, 245,
	250, 270, 280, 300, 300, 280, 270, 250,
	255, 300, 310, 315, 315, 310, 300, 255,
	260, 310, 320, 325, 325, 320, 310, 260,
	270, 320, 330, 330, 330, 330, 320, 270,
	270, 320, 330, 330, 330, 330, 320, 270,
	270, 320, 330, 330, 330, 330, 320, 270,
	270, 280, 285, 285, 285, 285, 280, 270,
}

// BISHOPS
var tableBishop = [NUM_SQUARES]Eval{
	250, 245, 250, 255, 255, 250, 245, 250,
	270, 300, 270, 290, 290, 270, 300, 270,
	285, 290, 295, 300, 300, 295, 290, 285,
	295, 305, 310, 310, 310, 310, 305, 295,
	305, 315, 320, 320, 320, 320, 315, 305,
	310, 320, 325, 325, 325, 325, 320, 310,
	300, 310, 320, 320, 320, 320, 310, 300,
	285, 285, 285, 285, 285, 285, 285, 285,
}

// ROOKS
var tableRook = [NUM_SQUARES]Eval{
	440, 460, 500, 515, 515, 500, 460, 440,
	440, 460, 500, 515, 515, 500, 460, 440,
	440, 460, 500, 515, 515, 500, 460, 440,
	440, 460, 500, 515, 515, 500, 460, 440,
	440, 460, 500, 515, 515, 500, 460, 440,
	440, 460, 500, 515, 515, 500, 460, 440,
	470, 490, 530, 545, 545, 530, 490, 470,
	460, 480, 520, 535, 535, 520, 480, 460,
}

// QUEENS
var tableQueen = [NUM_SQUARES]Eval{
	800, 825, 860, 880, 880, 860, 825, 800,
	825, 860, 880, 900, 900, 880, 860, 825,
	860, 880, 900, 925, 925, 900, 880, 860,
	880, 900, 925, 950, 950, 925, 900, 880,
	880, 900, 925, 950, 950, 925, 900, 880,
	860, 880, 900, 925, 925, 900, 880, 860,
	825, 860, 880, 900, 900, 880, 860, 825,
	800, 825, 860, 880, 880, 860, 825, 800,
}

// KINGS: Safety
var tableKingSafety = [NUM_SQUARES]Eval{
	350, 350, 300, 275, 275, 300, 350, 350,
	325, 325, 290, 290, 290, 290, 325, 325,
	295, 290, 250, 250, 250, 250, 290, 295,
	290, 250, 200, 200, 200, 200, 250, 290,
	250, 200, 150, 150, 150, 150, 200, 250,
	200, 150, 125, 125, 125, 125, 150, 200,
	150, 125, 100, 100, 100, 100, 125, 150,
	125, 100, 100, 100, 100, 100, 100, 125,
}

// KINGS: Active
var tableKingActive = [NUM_SQUARES]Eval{
	-50, -30, -30, -30, -30, -30, -30, -50,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-50, -40, -30, -20, -20, -30, -40, -50,
}

// KINGS: Survival
var tableKingSurvival = [NUM_SQUARES]Eval{
	-90, -90, -90, -90, -90, -90, -90, -90,
	-90, -50, -50, -50, -50, -50, -50, -90,
	-90, -50, 20, 20, 20, 20, -50, -90,
	-90, -50, 20, 50, 50, 20, -50, -90,
	-90, -50, 20, 50, 50, 20, -50, -90,
	-90, -50, 20, 20, 20, 20, -50, -90,
	-90, -50, -50, -50, -50, -50, -50, -90,
	-90, -90, -90, -90, -90, -90, -90, -90,
}

// Takes a white table and returns the black equivalent (flipped vertically)
func flipPST(table [NUM_SQUARES]Eval) [NUM_SQUARES]Eval {
	var newTable [NUM_SQUARES]Eval
	for i := 0; i < 64; i++ {
		// XOR 56 flips the rank (0->7, 1->6, etc)
		newTable[i^56] = table[i]
	}
	return newTable
}

// Master PST
var PST [NUM_GAME_STATES][NUM_COLORS][NUM_PIECES][NUM_SQUARES]Eval

func initPST() {
	// Define White's Tables

	// OPENING
	PST[OPENING][WHITE][PAWN] = tablePawn
	PST[OPENING][WHITE][KNIGHT] = tableKnight
	PST[OPENING][WHITE][BISHOP] = tableBishop
	PST[OPENING][WHITE][ROOK] = tableRook
	PST[OPENING][WHITE][QUEEN] = tableQueen
	PST[OPENING][WHITE][KING] = tableKingSafety

	// MIDDLEGAME
	PST[MIDDLEGAME][WHITE][PAWN] = tablePawn
	PST[MIDDLEGAME][WHITE][KNIGHT] = tableKnight
	PST[MIDDLEGAME][WHITE][BISHOP] = tableBishop
	PST[MIDDLEGAME][WHITE][ROOK] = tableRook
	PST[MIDDLEGAME][WHITE][QUEEN] = tableQueen
	PST[MIDDLEGAME][WHITE][KING] = tableKingSafety

	// ENDGAME
	PST[ENDGAME][WHITE][PAWN] = tablePawn
	PST[ENDGAME][WHITE][KNIGHT] = tableKnight
	PST[ENDGAME][WHITE][BISHOP] = tableBishop
	PST[ENDGAME][WHITE][ROOK] = tableRook
	PST[ENDGAME][WHITE][QUEEN] = tableQueen
	PST[ENDGAME][WHITE][KING] = tableKingActive

	// MATING
	PST[MATING][WHITE][PAWN] = tablePawn
	PST[MATING][WHITE][KNIGHT] = tableKnight
	PST[MATING][WHITE][BISHOP] = tableBishop
	PST[MATING][WHITE][ROOK] = tableRook
	PST[MATING][WHITE][QUEEN] = tableQueen
	PST[MATING][WHITE][KING] = tableKingActive

	// BEING MATED
	PST[BEING_MATED][WHITE][PAWN] = tablePawn
	PST[BEING_MATED][WHITE][KNIGHT] = tableKnight
	PST[BEING_MATED][WHITE][BISHOP] = tableBishop
	PST[BEING_MATED][WHITE][ROOK] = tableRook
	PST[BEING_MATED][WHITE][QUEEN] = tableQueen
	PST[BEING_MATED][WHITE][KING] = tableKingSurvival

	// Generate Black's tables by flipping White's
	for state := range int(NUM_GAME_STATES) {
		for piece := range int(NUM_PIECES) {
			PST[state][BLACK][piece] = flipPST(PST[state][WHITE][piece])
		}
	}
}
