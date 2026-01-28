package engine

/*
This file contains all the code related to searching
*/

// Simple negamax function for checking Perft, not actually used in searching and evalution, just testing
type PerftNegamaxResult struct {
	nodes int
}

func (b *Board) perftNegamax(depth uint8, moveStack [][]Move) PerftNegamaxResult {

	// Reaching a terminal condition
	if depth == 0 {
		return PerftNegamaxResult{
			nodes: 1,
		}
	}

	// Generate the pseudo legal moves to play, populating this depths move in the movestack
	nodes := 0
	moves := moveStack[depth]
	numberOfMoves := b.generatePseudoLegalMoves(moves)

	for _, move := range moves[:numberOfMoves] {
		unmake, isLegal := b.makeMove(move)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}
		result := b.perftNegamax(depth-1, moveStack)
		nodes += result.nodes
		b.unMakeMove(unmake)
	}

	return PerftNegamaxResult{
		nodes: nodes,
	}
}
