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

func (b *Board) rootSearch(depth uint8, multithread bool) RootSearchResult {

	// Validate depth is reasonable
	if depth == 0 {
		depth = 1
	} else if depth > 10 {
		depth = 10
	}

	// Allocate the moveStack
	moveStack := make([][]Move, MAX_PLY)
	for i := range moveStack {
		moveStack[i] = make([]Move, 256)
	}

	// Setting up the killer moves
	var killers Killers

	// Setting up the history cutoff heuristic
	var cutoffHistory CutoffHeuristic

	// Setup the search
	nodes := 1
	bestEval := MIN_EVAL
	alpha := MIN_EVAL
	beta := MAX_EVAL
	ply := uint8(0)

	// Check the TT table
	// This is not to prevent the entire root search, but to help move ordering
	var ttEntry *TTEntry = nil

	// Compute the TT key
	key := b.Zobrist & (TT_SIZE - 1)

	// Get the entry and return hit
	entry := &TT[key]
	if entry.zobrist == b.Zobrist {
		ttEntry = entry
	}

	// Generate the pseudo legal moves to play, populating this depths move in the movestack
	moves := moveStack[ply]
	numberOfMoves := b.generatePseudoLegalMovesWithOrdering(moves, ttEntry, nil, nil, nil)
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
		result := b.abnegamax(ply+1, depth-1, -beta, -alpha, moveStack, &killers, &cutoffHistory)
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

func (b *Board) abnegamax(ply uint8, depth uint8, alpha, beta Eval, moveStack [][]Move, killers *Killers, cutoffHistory *CutoffHeuristic) SearchResult {

	// checking for 3-fold repition
	// if it is, the game is a draw
	if b.isThreeFold() {
		return SearchResult{
			nodes: 1,
			best: MoveEval{
				move: Move{},
				eval: Eval(0),
			},
		}
	}

	// saving original alpha for TT tables
	originalAlpha := alpha
	originalBeta := beta

	// Check the TT table
	var ttEntry *TTEntry = nil

	// Compute the TT key
	key := b.Zobrist & (TT_SIZE - 1)

	// Get the entry and return hit
	entry := &TT[key]
	if entry.zobrist == b.Zobrist {
		ttEntry = entry
	}

	// Check if tt was found and was depth of equal or greater
	if ttEntry != nil && ttEntry.depth >= depth {

		switch ttEntry.flag {
		case TT_EXACT:
			return SearchResult{
				nodes: 1,
				best: MoveEval{
					move: ttEntry.move,
					eval: ttEntry.eval,
				},
			}
		case TT_LOWER:
			if ttEntry.eval >= beta {
				return SearchResult{
					nodes: 1,
					best: MoveEval{
						move: ttEntry.move,
						eval: ttEntry.eval,
					},
				}
			}
			alpha = max(alpha, ttEntry.eval)

		case TT_UPPER:
			if ttEntry.eval <= alpha {
				return SearchResult{
					nodes: 1,
					best: MoveEval{
						move: ttEntry.move,
						eval: ttEntry.eval,
					},
				}
			}
			beta = min(beta, ttEntry.eval)
		}

		// Check if the ab-window closed
		if alpha >= beta {
			return SearchResult{
				nodes: 1,
				best: MoveEval{
					move: ttEntry.move,
					eval: ttEntry.eval,
				},
			}
		}
	}

	// If at base condition, quiescence search
	if depth == 0 {
		return b.quiescence(ply+1, alpha, beta, moveStack)
	}

	// Setup the search
	nodes := 1
	bestEval := MIN_EVAL
	bestMove := Move{}

	// Two ply killers are killer moves from the previous position for this color
	var twoPlyKillers *[2]Move
	if ply >= 2 {
		twoPlyKillers = &(*killers)[ply-1]
	}
	thisKillers := (*killers)[ply]

	// Generate the pseudo legal moves to play, populating this plys move in the movestack
	moves := moveStack[ply]
	numberOfMoves := b.generatePseudoLegalMovesWithOrdering(moves, ttEntry, &thisKillers, twoPlyKillers, cutoffHistory)
	legalMovesFound := false
	for i, move := range moves[:numberOfMoves] {

		// Make the move and see if it was legal
		unmake, isLegal := b.makeMove(move)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}

		// Using late move reduction
		// Speeds up search 10x, costs 0.80 points on the benchmark test
		betaSearch := beta
		reduction := uint8(0)
		if i > 10 && depth > 2 && move.code != MOVE_CODE_CAPTURE && move != thisKillers[0] && move != thisKillers[1] {
			reduction = 1
			betaSearch = alpha + 1

			if i > 20 && depth > 2 {
				reduction = 2
			}
		}

		// Search the new position and get the results
		legalMovesFound = true
		result := b.abnegamax(ply+1, depth-1-reduction, -betaSearch, -alpha, moveStack, killers, cutoffHistory)
		resultEval := -result.best.eval
		nodes += result.nodes

		// If the engine reduced and the engine exceeded alpha, the engine needs to research at a full depth
		if resultEval > alpha && resultEval < beta && reduction > 0 {
			result = b.abnegamax(ply+1, depth-1, -beta, -alpha, moveStack, killers, cutoffHistory)
		}

		b.unMakeMove(unmake)
		if resultEval > bestEval {
			bestEval = resultEval
			bestMove = move
			if resultEval > alpha {
				alpha = resultEval
			}
		}

		// Failed soft on beta-cutoff, exit the search
		if resultEval > beta {
			bestEval = resultEval
			bestMove = move

			// Update killers
			// Make sure it is not a capture
			if move.code != MOVE_CODE_CAPTURE && move.code != MOVE_CODE_EN_PASSANT {
				if killers[ply][0] != move {
					killers[ply][1] = killers[ply][0]
					killers[ply][0] = move
				}

				// Update history of cutoffs as well (if not capture)
				cutoffHistory[b.Turn][move.start][move.target] += int(depth) * int(depth)
			}
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

	// Only update TT if searching at a greater or equal depth than previous entry
	if ttEntry == nil || depth >= ttEntry.depth {
		var ttFlag uint8
		if bestEval <= originalAlpha {
			ttFlag = TT_UPPER
		} else if bestEval >= originalBeta {
			ttFlag = TT_LOWER
		} else {
			ttFlag = TT_EXACT
		}
		// Store the value in the TT table
		updateTT(b.Zobrist, bestEval, ttFlag, depth, bestMove)
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

	// First, evalute the stand pat score of the position, the evaluation before doing any more captures
	standPat := b.eval()
	bestEval := standPat
	if b.Turn == BLACK {
		bestEval *= -1
	}

	// If the stand pat failed over beta, return it
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

	moves := moveStack[ply]
	numberOfMoves := b.generatePseudoLegalMovesWithOrdering(moves, nil, nil, nil, nil)
	for _, move := range moves[:numberOfMoves] {

		// Make sure the move was a capture
		if move.code != MOVE_CODE_CAPTURE || move.code == MOVE_CODE_EN_PASSANT {
			continue
		}

		// Delta pruning
		// If the capture for free, plus stand pat and a margin does not exceed alpha, do not search
		if standPat+PIECE_VALUES[b.MailBox[move.target]]+DELTA_MARGIN <= alpha {
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
