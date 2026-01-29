package engine

/*
This file contains all the code related to searching
*/

// Root search is the starting search for the chess engine, before it goes into its alpha-beta-negamax
// Here certain setup steps can take place, like multi-threading, if needed outside the main recursion
// It also handles validating search safety, so a depth of like 100 isn't run on the engine
// It returns the move and evals for all the root moves, up to the caller to determine if it should sort and slice or not
type RootSearchResult struct {
	nodes int
	moves []MoveEval
}

func (b *Board) rootSearch(depth uint8, moveStack [][]Move, multithread bool) RootSearchResult {

	// Validate depth is reasonable
	if depth == 0 {
		depth = 1
	} else if depth > 10 {
		depth = 10
	}

	// Setup the search
	nodes := 1
	bestEval := MIN_EVAL
	alpha := MIN_EVAL
	beta := MAX_EVAL
	ply := uint8(0)

	// Generate the pseudo legal moves to play, populating this depths move in the movestack
	moves := moveStack[ply]
	numberOfMoves := b.generatePseudoLegalMoves(moves)
	results := make([]MoveEval, 0, numberOfMoves)
	legalMovesFound := false
	for _, move := range moves[:numberOfMoves] {

		// Make the move and see if it was legal
		unmake, isLegal := b.makeMove(move)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}

		// Search the new position and get the results
		legalMovesFound = true
		result := b.abnegamax(ply+1, depth-1, -beta, -alpha, moveStack)
		b.unMakeMove(unmake)
		resultEval := -result.best.eval
		results = append(results, MoveEval{
			eval: resultEval,
			move: move,
		})
		nodes += result.nodes
		if resultEval > bestEval {
			bestEval = resultEval
			if resultEval > alpha {
				alpha = resultEval
			}
		}

		// Failed soft on beta-cutoff, exit the search
		// todo: remove this and beta, as beta never gets updated in rootSearch and thus this will not happen
		if resultEval >= beta {
			break
		}
	}

	// Handle checkmate/stalemate
	if !legalMovesFound {
		// If not in check, then stalement, else MIN_EVAL is correct
		if !b.isInCheck(b.Turn) {
			bestEval = 0
		}
	}

	return RootSearchResult{
		nodes: nodes,
		moves: results,
	}
}

// abnegamax is the main recursive search for the engine (negamax with alpha beta pruning)
// todo: abnegamax should be time aware
type SearchResult struct {
	nodes int
	best  MoveEval
}

func (b *Board) abnegamax(ply uint8, depth uint8, alpha, beta Eval, moveStack [][]Move) SearchResult {

	// If at base condition, quiescence search
	if depth == 0 {
		return b.quiescence(ply, alpha, beta, moveStack)
	}

	// Setup the search
	nodes := 1
	bestEval := MIN_EVAL
	bestMove := Move{}

	// Generate the pseudo legal moves to play, populating this plys move in the movestack
	moves := moveStack[ply]
	numberOfMoves := b.generatePseudoLegalMoves(moves)
	legalMovesFound := false
	for _, move := range moves[:numberOfMoves] {

		// Make the move and see if it was legal
		unmake, isLegal := b.makeMove(move)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}

		// Search the new position and get the results
		legalMovesFound = true
		result := b.abnegamax(ply+1, depth-1, -beta, -alpha, moveStack)
		b.unMakeMove(unmake)
		resultEval := -result.best.eval
		nodes += result.nodes
		if resultEval > bestEval {
			bestEval = resultEval
			bestMove = move
			if resultEval > alpha {
				alpha = resultEval
			}
		}

		// Failed soft on beta-cutoff, exit the search
		if resultEval >= beta {
			break
		}
	}

	// Handle checkmate/stalemate
	if !legalMovesFound {
		// If not in check, then stalement
		if !b.isInCheck(b.Turn) {
			bestEval = 0
		} else {
			// If is in check, then take the MIN_EVAL and add the ply to it to prioritize faster mates
			bestEval += Eval(ply)
		}
	}

	return SearchResult{
		nodes: nodes,
		best: MoveEval{
			move: bestMove,
			eval: bestEval,
		},
	}
}

// quiescence is the final search for a "quiet" position the engine takes, after reaching the base condition of abnegamax
// A quiet position is one without any captures
// todo: should be upgraded to check for checks as well
func (b *Board) quiescence(ply uint8, alpha, beta Eval, moveStack [][]Move) SearchResult {

	// b.eval is absolute: negative is good for black and psotive is good for white
	// however, quiescenece and negamax need to return "context aware" evals
	// as such, we need to negate the eval if the board we are evaluating is black
	// negamax returns the score as relation to positive being good for the active play
	// so we must flip blacks eval sign
	// if no captures found, at end of quiescence search and should evaluate
	bestEval := b.eval()
	if b.Turn == BLACK {
		bestEval *= -1
	}

	// check if that caused a soft-beta cutoff
	if bestEval >= beta {
		return SearchResult{
			nodes: 1,
			best:  MoveEval{eval: bestEval},
		}
	}

	// update alpha if needed
	if bestEval > alpha {
		alpha = bestEval
	}

	// check ply, if it exceeds or equals MAX_PLY then just evalute
	// this is just a safety net against really weird conditions, very unlikely to happen
	nodes := 1
	bestMove := Move{}
	if ply >= MAX_PLY {
		return SearchResult{
			nodes: 1,
			best: MoveEval{
				move: bestMove,
				eval: bestEval,
			},
		}
	}

	// For now quiescence search will just evaluate
	moves := moveStack[ply]
	numberOfMoves := b.generatePseudoLegalMoves(moves)
	for _, move := range moves[:numberOfMoves] {

		// Make sure the move was a capture
		if move.code != MOVE_CODE_CAPTURE && move.code != MOVE_CODE_EN_PASSANT {
			continue
		}

		// Make the move and see if it was legal
		unmake, isLegal := b.makeMove(move)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}

		// Search the new position and get the results
		result := b.quiescence(ply+1, -beta, -alpha, moveStack)
		b.unMakeMove(unmake)
		resultEval := -result.best.eval
		nodes += result.nodes
		if resultEval > bestEval {
			bestEval = resultEval
			bestMove = move
			if resultEval > alpha {
				alpha = resultEval
			}
		}

		// Failed soft on beta-cutoff, exit the search
		if resultEval >= beta {
			break
		}
	}

	return SearchResult{
		nodes: nodes,
		best: MoveEval{
			move: bestMove,
			eval: bestEval,
		},
	}
}

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
