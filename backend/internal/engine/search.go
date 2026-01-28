package engine

/*
This file contains all the code related to searching
*/

// Negamax is the main recursive function used to search
type NegamaxResult struct {
	eval  Eval
	nodes int
}

func (b *Board) negamax(depth uint8, moveStack [][]Move) NegamaxResult {

	// Reaching a terminal condition
	if depth == 0 {
		return NegamaxResult{
			eval:  Eval(0),
			nodes: 1,
		}
	}

	eval := MIN_EVAL
	nodes := 0

	// Generate the pseudo legal moves to play, populating this depths move in the movestack
	// TODO: Sort these moves
	moves := moveStack[depth]
	numberOfMoves := b.generatePseudoLegalMoves(moves)

	for _, move := range moves[:numberOfMoves] {
		unmake, isLegal := b.makeMove(move)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}
		result := b.negamax(depth-1, moveStack)
		eval = max(eval, -result.eval)
		nodes += result.nodes
		b.unMakeMove(unmake)
	}

	return NegamaxResult{
		eval:  eval,
		nodes: nodes,
	}
}
