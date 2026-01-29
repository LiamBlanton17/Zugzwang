package engine

/*
This file holds the global Piece-Square Tables (psts) for the evaluation function
*/

/*
Referenece for the board in numbers. It is mirrored with relation to the tables below (index 0 is bottom left of board)

56 57 58 59 60 61 62 63
48 49 50 51 52 53 54 55
40 41 42 43 44 45 46 47
32 33 34 35 36 37 38 39
24 25 26 27 28 29 30 31
16 17 18 19 20 21 22 23
08 09 10 11 12 13 14 15
00 01 02 03 04 05 06 07
*/

// PAWNS
// Prioritize pawn advancements on the queenside and center
var pstPawnOpening = [NUM_SQUARES]Eval{
	000, 000, 000, 000, 000, 000, 000, 000,
	90, 90, 90, 80, 80, 120, 100, 100,
	95, 100, 100, 90, 90, 70, 100, 100,
	105, 105, 105, 130, 130, 85, 90, 95,
	110, 110, 110, 140, 140, 100, 105, 105,
	115, 120, 120, 145, 145, 115, 115, 115,
	160, 160, 160, 160, 160, 160, 160, 160,
	000, 000, 000, 000, 000, 000, 000, 000,
}

// Prioritize pawn advancements on the edges more than the center
// But still score storng for center pawns
var pstPawnEndgame = [NUM_SQUARES]Eval{
	000, 000, 000, 000, 000, 000, 000, 000,
	80, 80, 80, 65, 65, 100, 95, 95,
	90, 90, 90, 90, 90, 90, 100, 100,
	105, 105, 110, 130, 130, 100, 105, 105,
	115, 115, 125, 140, 140, 120, 115, 115,
	145, 145, 145, 145, 145, 140, 145, 145,
	185, 180, 170, 165, 165, 170, 180, 185,
	000, 000, 000, 000, 000, 000, 000, 000,
}

// KNIGHTS
// Priortize development, center and forward, and staying away from the edge
var pstKnightOpening = [NUM_SQUARES]Eval{
	200, 235, 245, 250, 250, 245, 235, 200,
	235, 245, 250, 295, 295, 250, 245, 235,
	245, 270, 315, 320, 320, 315, 270, 245,
	250, 320, 320, 325, 325, 320, 320, 250,
	255, 320, 325, 330, 330, 325, 320, 255,
	260, 325, 335, 335, 335, 335, 325, 260,
	255, 305, 315, 320, 320, 315, 305, 255,
	225, 245, 260, 260, 260, 260, 245, 225,
}

// Priortize centrality above all else
var pstKnightEndgame = [NUM_SQUARES]Eval{
	200, 225, 250, 275, 275, 250, 225, 200,
	225, 255, 280, 305, 305, 280, 255, 225,
	250, 280, 310, 320, 320, 310, 280, 250,
	275, 305, 320, 335, 325, 320, 305, 275,
	275, 305, 320, 335, 325, 320, 305, 275,
	250, 280, 310, 320, 320, 310, 280, 250,
	225, 255, 280, 305, 305, 280, 255, 225,
	200, 225, 250, 275, 275, 250, 225, 200,
}

// BISHOPS
// Priortize development to key squares (b2/g2/c4/f4/b5/g5)
var pstBishopOpening = [NUM_SQUARES]Eval{
	215, 220, 270, 265, 265, 270, 220, 215,
	280, 295, 275, 290, 290, 250, 305, 280,
	290, 290, 290, 295, 295, 290, 290, 290,
	300, 305, 315, 320, 320, 315, 305, 300,
	305, 315, 320, 320, 320, 320, 315, 305,
	300, 305, 310, 310, 310, 310, 305, 300,
	285, 295, 300, 300, 300, 300, 295, 285,
	270, 270, 275, 275, 275, 275, 270, 270,
}

// Priortize centrality above all else
var pstBishopEndgame = [NUM_SQUARES]Eval{
	215, 245, 265, 285, 285, 265, 245, 215,
	245, 270, 295, 320, 320, 295, 270, 245,
	265, 295, 325, 340, 340, 325, 295, 250,
	285, 320, 340, 355, 355, 340, 320, 285,
	285, 320, 340, 355, 355, 340, 320, 285,
	265, 295, 325, 340, 340, 325, 295, 250,
	245, 270, 295, 320, 320, 295, 270, 245,
	215, 245, 265, 285, 285, 265, 245, 215,
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
