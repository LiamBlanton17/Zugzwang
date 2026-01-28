package engine

/*
This file holds the global Piece-Square Tables (psts) for the evaluation function
*/

// Base tables

// PAWNS
var tablePawn = [NUM_SQUARES]Eval{
	0, 0, 0, 0, 0, 0, 0, 0, // Rank 1
	0, 10, 20, -50, -50, 20, 10, 0, // Rank 2
	5, 15, -50, 25, 25, 15, 15, 5, // Rank 3
	0, 0, -10, 50, 50, 0, 0, 0, // Rank 4
	10, 10, 20, 50, 50, 20, 10, 10, // Rank 5
	20, 20, 45, 70, 70, 45, 20, 20, // Rank 6
	75, 75, 75, 75, 75, 75, 75, 75, // Rank 7
	0, 0, 0, 0, 0, 0, 0, 0, // Rank 8
}

// KNIGHTS
var tableKnight = [NUM_SQUARES]Eval{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

// BISHOPS
var tableBishop = [NUM_SQUARES]Eval{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

// ROOKS
var tableRook = [NUM_SQUARES]Eval{
	0, 0, 0, 5, 5, 0, 0, 0,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	5, 10, 10, 10, 10, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// QUEENS
var tableQueen = [NUM_SQUARES]Eval{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-10, 5, 5, 5, 5, 5, 0, -10,
	0, 0, 5, 5, 5, 5, 0, -5,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}

// KINGS: Safety
var tableKingSafety = [NUM_SQUARES]Eval{
	20, 30, 10, 0, 0, 10, 30, 20,
	20, 20, 0, 0, 0, 0, 20, 20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
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
