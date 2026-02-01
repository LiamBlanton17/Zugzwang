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
	100, 100, 95, 90, 90, 105, 100, 100,
	95, 95, 100, 100, 100, 80, 95, 95,
	95, 100, 105, 110, 110, 85, 90, 90,
	95, 100, 110, 115, 115, 90, 90, 90,
	100, 105, 115, 120, 120, 100, 95, 95,
	105, 110, 120, 125, 125, 110, 100, 100,
	000, 000, 000, 000, 000, 000, 000, 000,
}

// Prioritize pawn advancements on the edges more than the center
// But still score storng for center pawns
var pstPawnEndgame = [NUM_SQUARES]Eval{
	000, 000, 000, 000, 000, 000, 000, 000,
	85, 85, 85, 85, 85, 85, 85, 85,
	95, 95, 95, 95, 95, 95, 95, 95,
	100, 100, 100, 100, 100, 100, 100, 100,
	105, 105, 105, 105, 105, 105, 105, 105,
	115, 115, 115, 115, 115, 115, 115, 115,
	135, 125, 125, 125, 125, 125, 130, 135,
	000, 000, 000, 000, 000, 000, 000, 000,
}

// KNIGHTS
// Priortize development, center and forward, and staying away from the edge
var pstKnightOpening = [NUM_SQUARES]Eval{
	280, 285, 285, 290, 290, 285, 285, 280,
	290, 290, 295, 300, 300, 295, 290, 290,
	295, 305, 310, 310, 310, 310, 305, 295,
	295, 310, 315, 315, 315, 315, 310, 295,
	295, 315, 315, 315, 315, 315, 315, 295,
	290, 305, 305, 305, 305, 305, 305, 290,
	285, 295, 300, 300, 300, 300, 295, 285,
	280, 285, 295, 295, 295, 295, 285, 280,
}

// Priortize centrality above all else
var pstKnightEndgame = [NUM_SQUARES]Eval{
	270, 280, 290, 300, 300, 290, 280, 270,
	280, 290, 300, 310, 310, 300, 290, 280,
	290, 300, 310, 315, 315, 310, 300, 290,
	300, 310, 315, 320, 320, 315, 310, 300,
	300, 310, 315, 320, 320, 315, 310, 300,
	290, 300, 310, 315, 315, 310, 300, 290,
	280, 290, 300, 310, 310, 300, 290, 280,
	270, 280, 290, 300, 300, 290, 280, 270,
}

// BISHOPS
// Priortize development to key squares (b2/g2/c4/f4/b5/g5)
var pstBishopOpening = [NUM_SQUARES]Eval{
	275, 285, 290, 290, 290, 290, 285, 275,
	285, 310, 285, 295, 295, 285, 310, 285,
	295, 295, 300, 305, 305, 300, 295, 295,
	300, 300, 305, 310, 310, 305, 300, 300,
	305, 305, 315, 315, 315, 315, 305, 305,
	305, 315, 315, 315, 315, 315, 315, 305,
	300, 310, 315, 315, 315, 315, 310, 300,
	295, 300, 305, 305, 305, 305, 300, 295,
}

// Priortize centrality above all else
// Also bishop overall value increases
var pstBishopEndgame = [NUM_SQUARES]Eval{
	275, 285, 295, 305, 305, 295, 285, 275,
	285, 295, 305, 315, 315, 305, 295, 285,
	295, 305, 315, 320, 320, 315, 305, 295,
	305, 315, 320, 325, 325, 320, 315, 305,
	305, 315, 320, 325, 325, 320, 315, 305,
	295, 305, 315, 320, 320, 315, 305, 295,
	285, 295, 305, 315, 315, 305, 295, 285,
	275, 285, 295, 305, 305, 295, 285, 275,
}

// ROOKS
// Priortize central movement, don't encourage rook lifts, but 7/8th rank do get bonuses
var pstRookOpening = [NUM_SQUARES]Eval{
	490, 495, 500, 505, 505, 500, 495, 490,
	495, 500, 505, 510, 510, 505, 500, 495,
	500, 505, 505, 510, 510, 505, 505, 500,
	500, 505, 510, 515, 515, 510, 505, 500,
	500, 505, 510, 515, 515, 510, 505, 500,
	505, 510, 515, 520, 520, 515, 510, 505,
	520, 525, 530, 535, 535, 530, 525, 520,
	510, 515, 520, 525, 525, 520, 515, 510,
}

// Slightly priortize centrality, with bonuses for 7/8th rank
// Also rooks overall value increases
var pstRookEndgame = [NUM_SQUARES]Eval{
	500, 505, 510, 515, 515, 510, 505, 500,
	510, 520, 525, 530, 530, 525, 520, 510,
	515, 525, 535, 540, 540, 535, 525, 515,
	520, 530, 540, 550, 550, 540, 530, 520,
	520, 530, 540, 550, 550, 540, 530, 520,
	515, 525, 535, 540, 540, 535, 525, 515,
	530, 540, 550, 560, 560, 550, 540, 530,
	510, 520, 530, 540, 540, 530, 520, 510,
}

// QUEENS
// Prioritize centrality, but not overzelous development
var pstQueenOpening = [NUM_SQUARES]Eval{
	840, 860, 870, 880, 880, 870, 860, 840,
	860, 880, 890, 895, 895, 890, 880, 860,
	870, 890, 900, 905, 905, 900, 890, 870,
	875, 895, 905, 910, 910, 905, 895, 875,
	875, 895, 905, 910, 910, 905, 895, 875,
	870, 890, 900, 905, 905, 900, 890, 870,
	860, 880, 890, 895, 895, 890, 880, 860,
	840, 860, 870, 880, 880, 870, 860, 840,
}

// Priortize centrality above all else
var pstQueenEndgame = [NUM_SQUARES]Eval{
	860, 880, 900, 915, 915, 900, 880, 860,
	880, 905, 920, 930, 930, 920, 905, 880,
	900, 920, 935, 945, 945, 935, 920, 900,
	915, 930, 945, 955, 955, 945, 930, 915,
	915, 930, 945, 955, 955, 945, 930, 915,
	900, 920, 935, 945, 945, 935, 920, 900,
	880, 905, 920, 930, 930, 920, 905, 880,
	860, 880, 900, 915, 915, 900, 880, 860,
}

// KINGS
// Prioritize safety above all else
var pstKingOpening = [NUM_SQUARES]Eval{
	500, 515, 500, 460, 460, 480, 515, 500,
	485, 490, 470, 445, 445, 455, 470, 485,
	460, 455, 440, 425, 425, 435, 455, 460,
	440, 435, 420, 405, 405, 415, 435, 440,
	420, 415, 400, 385, 385, 395, 415, 420,
	400, 395, 380, 365, 365, 375, 395, 400,
	380, 375, 360, 350, 350, 360, 375, 380,
	370, 370, 360, 350, 350, 360, 370, 370,
}

// Prioritize center above all else
var pstKingEndgame = [NUM_SQUARES]Eval{
	340, 370, 400, 420, 420, 400, 370, 340,
	370, 400, 430, 450, 450, 430, 400, 370,
	400, 430, 460, 480, 480, 460, 430, 400,
	420, 450, 480, 500, 500, 480, 450, 420,
	420, 450, 480, 500, 500, 480, 450, 420,
	400, 430, 460, 480, 480, 460, 430, 400,
	370, 400, 430, 450, 450, 430, 400, 370,
	340, 370, 400, 420, 420, 400, 370, 340,
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
var PST [2][NUM_COLORS][NUM_PIECES][NUM_SQUARES]Eval

// Init pst tables, to be called on engine startup
func initPST() {
	// Define White's Tables
	// OPENING
	PST[OPENING][WHITE][PAWN] = pstPawnOpening
	PST[OPENING][WHITE][KNIGHT] = pstKnightOpening
	PST[OPENING][WHITE][BISHOP] = pstBishopOpening
	PST[OPENING][WHITE][ROOK] = pstRookOpening
	PST[OPENING][WHITE][QUEEN] = pstQueenOpening
	PST[OPENING][WHITE][KING] = pstKingOpening

	// ENDGAME
	PST[ENDGAME][WHITE][PAWN] = pstPawnEndgame
	PST[ENDGAME][WHITE][KNIGHT] = pstKnightEndgame
	PST[ENDGAME][WHITE][BISHOP] = pstBishopEndgame
	PST[ENDGAME][WHITE][ROOK] = pstRookEndgame
	PST[ENDGAME][WHITE][QUEEN] = pstQueenEndgame
	PST[ENDGAME][WHITE][KING] = pstKingEndgame

	// Generate Black's tables by flipping White's
	for state := range 2 {
		for piece := range int(NUM_PIECES) {
			PST[state][BLACK][piece] = flipPST(PST[state][WHITE][piece])
		}
	}
}
