package engine

import "time"

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

/*
Deprecated (at least I think lol, don't use)
Evaluate is the standard function to evalute a position, to be used by the API package to utilize the engine.
*/
type EvaluateResponse struct {
	MoveEvals []MoveEval
	duration  time.Duration
	nodes     int32
}

func Evalute(position FEN, history []FEN, numberOfMoves int) (*EvaluateResponse, error) {
	// Time the function from start to button, including building the board and history
	// This is done as this provides a more accurate evalution of how fast the engine is
	start := time.Now()

	// Build the board from the position
	// This can fail if position is not a valid FEN string
	board, err := position.ToBoard(history)
	if err != nil {
		return nil, err
	}

	// Search the board and get the results
	results := board.search(numberOfMoves)

	// Stop the time
	end := time.Now()

	return &EvaluateResponse{
		MoveEvals: results.MoveEvals,
		duration:  end.Sub(start),
		nodes:     results.Nodes,
	}, nil
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
