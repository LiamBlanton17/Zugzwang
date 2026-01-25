package engine

/*
This file contains functionality related to setup, searching and evalution of a board.
*/

func (f FEN) toZobrist() ZobristHash {

}

func (position FEN) toBoard() *Board {
	board := Board{
		Zobrist: position.toZobrist(),
	}

	return &board
}

func buildGameHistory(history []FEN) *GameHistory {
	var gameHistory GameHistory

	for _, h := range history {
		gameHistory = append(gameHistory, h.toZobrist())
	}

	return &gameHistory
}

type BoardSearchResults struct {
	Nodes     int32
	MoveEvals []MoveEval
}

func (b *Board) search(history *GameHistory, numberOfMoves int) BoardSearchResults {
	return BoardSearchResults{}
}

func (b *Board) generateMoves() []Move {
	var moves []Move

	return moves
}
