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
// Also bishop overall value increases
var pstBishopEndgame = [NUM_SQUARES]Eval{
	220, 245, 265, 285, 285, 265, 245, 220,
	245, 270, 295, 320, 320, 295, 270, 245,
	265, 295, 325, 340, 340, 325, 295, 250,
	285, 320, 340, 355, 355, 340, 320, 285,
	285, 320, 340, 355, 355, 340, 320, 285,
	265, 295, 325, 340, 340, 325, 295, 250,
	245, 270, 295, 320, 320, 295, 270, 245,
	220, 245, 265, 285, 285, 265, 245, 220,
}

// ROOKS
// Priortize central movement, don't encourage rook lifts, but 7/8th rank do get bonuses
var pstRookOpening = [NUM_SQUARES]Eval{
	470, 490, 510, 520, 520, 515, 480, 460,
	440, 440, 500, 505, 505, 500, 440, 440,
	440, 440, 500, 505, 505, 500, 440, 440,
	440, 440, 500, 505, 505, 500, 440, 440,
	440, 440, 500, 505, 505, 500, 440, 440,
	440, 440, 500, 505, 505, 500, 440, 440,
	490, 490, 510, 515, 515, 510, 490, 490,
	490, 490, 510, 515, 515, 510, 490, 490,
}

// Slightly priortize centrality, with bonuses for 7/8th rank
// Also rooks overall value increases
var pstRookEndgame = [NUM_SQUARES]Eval{
	490, 500, 510, 520, 520, 510, 500, 490,
	500, 515, 530, 540, 540, 530, 515, 500,
	510, 530, 540, 550, 550, 540, 530, 510,
	520, 540, 550, 555, 555, 550, 540, 520,
	520, 540, 550, 555, 555, 550, 540, 520,
	510, 530, 540, 550, 550, 540, 530, 510,
	535, 550, 565, 570, 570, 565, 550, 535,
	515, 525, 535, 540, 540, 535, 525, 515,
}

// QUEENS
// Prioritize centrality, but not overzelous development
var pstQueenOpening = [NUM_SQUARES]Eval{
	850, 880, 900, 900, 880, 880, 880, 850,
	855, 885, 905, 905, 900, 885, 885, 855,
	870, 900, 910, 910, 910, 910, 900, 870,
	875, 905, 915, 915, 915, 915, 905, 875,
	875, 905, 915, 915, 915, 915, 905, 875,
	875, 905, 915, 915, 915, 915, 905, 875,
	865, 895, 900, 900, 900, 900, 895, 865,
	850, 880, 885, 885, 885, 885, 880, 850,
}

// Priortize centrality above all else
var pstQueenEndgame = [NUM_SQUARES]Eval{
	865, 885, 890, 915, 915, 890, 885, 865,
	885, 905, 920, 930, 930, 920, 905, 885,
	890, 920, 935, 945, 945, 935, 920, 890,
	915, 930, 945, 950, 950, 945, 930, 915,
	915, 930, 945, 950, 950, 945, 930, 915,
	890, 920, 935, 945, 945, 935, 920, 890,
	885, 905, 920, 930, 930, 920, 905, 885,
	865, 885, 890, 915, 915, 890, 885, 865,
}

// KINGS
// Prioritize safety above all else
var pstKingOpening = [NUM_SQUARES]Eval{
	525, 550, 525, 470, 470, 490, 550, 525,
	500, 500, 490, 465, 465, 465, 490, 500,
	460, 460, 450, 450, 450, 450, 460, 460,
	440, 440, 420, 420, 420, 420, 440, 440,
	420, 420, 390, 390, 390, 390, 420, 420,
	400, 400, 370, 370, 370, 370, 400, 400,
	375, 375, 350, 350, 350, 350, 375, 375,
	350, 350, 350, 350, 350, 350, 350, 350,
}

// Prioritize center above all else
var pstKingEndgame = [NUM_SQUARES]Eval{
	250, 350, 400, 450, 450, 400, 350, 250,
	350, 400, 450, 500, 500, 450, 400, 350,
	400, 450, 500, 525, 525, 500, 450, 400,
	450, 500, 525, 550, 550, 525, 500, 450,
	450, 500, 525, 550, 550, 525, 500, 450,
	400, 450, 500, 525, 525, 500, 450, 400,
	350, 400, 450, 500, 500, 450, 400, 350,
	250, 350, 400, 450, 450, 400, 350, 250,
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
