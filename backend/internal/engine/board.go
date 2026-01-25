package engine

/*
This file contains functionality related to setup, searching and evalution of a board.
*/

func buildBoard(position FEN) *Board {
	return nil
}

func buildGameHistory(history []FEN) *GameHistory {
	return nil
}

type BoardSearchResults struct {
	Nodes     int32
	MoveEvals []MoveEval
}

func (b *Board) search(history *GameHistory, numberOfMoves int) BoardSearchResults {
	return BoardSearchResults{}
}
