package engine

/*
This file contains the API to use the engine
*/

// This is pretty much just an alias for RootSearch() in search.go, so maybe just use that
func (b *Board) Search(depth uint8, multithread bool) RootSearchResult {
	return b.RootSearch(depth, multithread)
}

// Get the legal moves for the board, this is just an alias for GenerateLegalMoves in board.go, so maybe can just just use that function
func (b *Board) Moves() []Move {
	return b.GenerateLegalMoves()
}

// This function is called by the HTTP api
func GetBestMoves(pgn string, numMoves int) []string {

	// Parse the PGN string into a board

	// Search the board

	// Return the best move(s) as a string
	return []string{"e2e4", "d2d4"}
}

/*
InitEngine should be called once at startup.
This setups globals like TT tables, Zobrist keys, and pregenerated moves
*/
func InitEngine() {

	// Setup global zobrist hashing
	initZobrist()

	// Setup move lookup tables
	initKnightMoves()
	initKingMoves()
	initMagicRook()
	initMagicBishop()

	// Setup PSTs
	initPST()

	// Setup eval
	initEval()

	// Setup TT
	initTT()
}
