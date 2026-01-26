package engine

/*
This file contains all the logic related to move generation, aside from the method on the board object
which actually creates the moves. That function utilizes some functions here to generate moves for a given board
*/

/*
Referenece for the board in numbers
56 57 58 59 60 61 62 63
48 49 50 51 52 53 54 55
40 41 42 43 44 45 46 47
32 33 34 35 36 37 38 39
24 25 26 27 28 29 30 31
16 17 18 19 20 21 22 23
08 09 10 11 12 13 14 15
00 01 02 03 04 05 06 07

With the starting position
r n b q k b n r
p p p p p p p p
- - - - - - - -
- - - - - - - -
- - - - - - - -
- - - - - - - -
P P P P P P P P
R N B Q K B N R

With human readable notation
a8 b8 c8 d8 e8 f8 g8 h8
a7 b7 c7 d7 e7 f7 g7 h7
a6 b6 c6 d6 e6 f6 g6 h6
a5 b5 c5 d5 e5 f5 g5 h5
a4 b4 c4 d4 e4 f4 g4 h4
a3 b3 c3 d3 e3 f3 g3 h3
a2 b2 c2 d2 e2 f2 g2 h2
a1 b1 c1 d1 e1 f1 g1 h1
*/

// Used by board.generateMoves() to get the pseudo-legal pawn moves
// Does include captures
// Does include promotion
// Does include enpassent
func getPawnMoves(pawns, enemyPieces, occupancy BitBoard, EPS Square, color Color, moves []Move, moveIdx int) int {

	for pawns > 0 {
		start := pawns.popSquare()

		if color == WHITE { // WHITE pawn moves
			oneSq := start + 8
			twoSq := start + 16

			// Single Push (Not promotions, those are handled at the end of the code)
			if (occupancy & (1 << oneSq)) == 0 {
				if oneSq <= 55 {
					moveIdx = addMove(moves, start, oneSq, MOVE_CODE_NONE, false, moveIdx)

					// Double Push (on second rank)
					if start >= 8 && start <= 15 && (occupancy&(1<<twoSq)) == 0 {
						moveIdx = addMove(moves, start, twoSq, MOVE_CODE_DOUBLE_PAWN_PUSH, false, moveIdx)
					}
				}
			}

			// Capture targets
			capLeft := start + 7
			capRight := start + 9

			// Prevent wrapping around board by making sure start column is not the A column (col 0)
			canCapLeft := start%8 > 0 && ((enemyPieces&(1<<capLeft)) != 0 || capLeft == EPS)

			// Prevent wrapping around board by making sure start column is not the H column (col 7)
			canCapRight := start%8 < 7 && ((enemyPieces&(1<<capRight)) != 0 || capRight == EPS)

			// Captures (Not promotion)
			if oneSq <= 55 {
				if canCapLeft {
					code := MOVE_CODE_CAPTURE
					if capLeft == EPS {
						code = MOVE_CODE_EN_PASSANT
					}
					moveIdx = addMove(moves, start, capLeft, code, false, moveIdx)
				}
				if canCapRight {
					code := MOVE_CODE_CAPTURE
					if capRight == EPS {
						code = MOVE_CODE_EN_PASSANT
					}
					moveIdx = addMove(moves, start, capRight, code, false, moveIdx)
				}
			}

			// Handle all promotions (Push or Capture landing on the last rank)
			if oneSq > 55 {
				// Push Promotion
				if (occupancy & (1 << oneSq)) == 0 {
					moveIdx = addMove(moves, start, oneSq, MOVE_CODE_PROMOTION, true, moveIdx)
				}
				// Capture Left Promotion
				if canCapLeft {
					moveIdx = addMove(moves, start, capLeft, MOVE_CODE_CAPTURE, true, moveIdx)
				}
				// Capture Right Promotion
				if canCapRight {
					moveIdx = addMove(moves, start, capRight, MOVE_CODE_CAPTURE, true, moveIdx)
				}
			}

		} else { // BLACK pawn moves
			oneSq := start - 8
			twoSq := start - 16

			// Single Push (Not promotions, those are handled at the end of the code)
			if (occupancy & (1 << oneSq)) == 0 {
				if oneSq >= 8 {
					moveIdx = addMove(moves, start, oneSq, MOVE_CODE_NONE, false, moveIdx)
					if start >= 48 && start <= 55 && (occupancy&(1<<twoSq)) == 0 {
						moveIdx = addMove(moves, start, twoSq, MOVE_CODE_DOUBLE_PAWN_PUSH, false, moveIdx)
					}
				}
			}

			// Capture targets
			capRight := start - 7
			capLeft := start - 9

			// Prevent wrapping around board by making sure start column is not the A column (col 0)
			canCapRight := start%8 < 7 && ((enemyPieces&(1<<capRight)) != 0 || capRight == EPS)

			// Prevent wrapping around board by making sure start column is not the H column (col 7)
			canCapLeft := start%8 > 0 && ((enemyPieces&(1<<capLeft)) != 0 || capLeft == EPS)

			// Captures (Not promotion)
			if oneSq >= 8 {
				if canCapLeft {
					code := MOVE_CODE_CAPTURE
					if capLeft == EPS {
						code = MOVE_CODE_EN_PASSANT
					}
					moveIdx = addMove(moves, start, capLeft, code, false, moveIdx)
				}
				if canCapRight {
					code := MOVE_CODE_CAPTURE
					if capRight == EPS {
						code = MOVE_CODE_EN_PASSANT
					}
					moveIdx = addMove(moves, start, capRight, code, false, moveIdx)
				}
			}

			// Handle all promotions (Push or Capture landing on the last rank)
			if oneSq < 8 {
				// Push Promotion
				if (occupancy & (1 << oneSq)) == 0 {
					moveIdx = addMove(moves, start, oneSq, MOVE_CODE_PROMOTION, true, moveIdx)
				}
				// Capture Left Promotion
				if canCapLeft {
					moveIdx = addMove(moves, start, capLeft, MOVE_CODE_CAPTURE, true, moveIdx)
				}
				// Capture Right Promotion
				if canCapRight {
					moveIdx = addMove(moves, start, capRight, MOVE_CODE_CAPTURE, true, moveIdx)
				}
			}
		}
	}

	return moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal king moves
// This does not include castling
func getKingMoves(king, friendlyPieces BitBoard, moves []Move, moveIdx int) int {

	for king > 0 {
		start := king.popSquare()
		targets := KING_MOVES[start] &^ friendlyPieces

		for targets > 0 {
			moveIdx = addMove(moves, start, targets.popSquare(), MOVE_CODE_NONE, false, moveIdx)
		}
	}

	return moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal knight moves
func getKnightMoves(knights, friendlyPieces BitBoard, moves []Move, moveIdx int) int {

	for knights > 0 {
		start := knights.popSquare()
		targets := KNIGHT_MOVES[start] &^ friendlyPieces

		for targets > 0 {
			moveIdx = addMove(moves, start, targets.popSquare(), MOVE_CODE_NONE, false, moveIdx)
		}
	}

	return moveIdx
}

// Helper function to add a move and increment the counter
func addMove(moves []Move, start, target Square, code uint8, is_promotion bool, moveIdx int) int {
	if is_promotion {
		for _, piece := range []Piece{KNIGHT, BISHOP, ROOK, QUEEN} {
			moves[moveIdx].start = start
			moves[moveIdx].target = target
			moves[moveIdx].promotion = piece
			moves[moveIdx].code = MOVE_CODE_PROMOTION
			moveIdx++
		}
	} else {
		moves[moveIdx].start = start
		moves[moveIdx].target = target
		moves[moveIdx].code = code
		moves[moveIdx].promotion = NO_PIECE
		moveIdx++
	}
	return moveIdx
}

// Global King lookup table
// Knight move generation logic is simple, can only move 1 direction each way
// Castling is handled elsewhere
// During move generation at run time, AND NOT the bitboard with the occupancy of friendly pieces to avoid capturing yourself
var KING_MOVES [64]BitBoard

func initKingMoves() {
	// UP-RIGHT +9
	// UP +8
	// UP-LEFT +7
	// RIGHT +1
	// DOWN-RIGHT -7
	// DONW -8
	// DOWN-LEFT -9
	// LEFT -1
	directions := []int{9, 8, 7, 1, -7, -8, -9, -1}
	for i := range NUM_SQUARES {
		for _, d := range directions {
			t := i + d // Target square

			// Check out of bounds (top/bottom)
			if t < 0 || t > 63 {
				continue
			}

			// Check out of bounds (crossing left/right)
			// Make sure the column shift was only 2 or less columns
			colShift := (t % 8) - (i % 8)
			if colShift < -1 || colShift > 1 {
				continue
			}

			KING_MOVES[i] |= Square(t).bitBoardPosition()
		}
	}
}

// Global Knight lookup table
// Knight move generation logic is simple, just map 64 squares to a set of moves
// During move generation at run time, AND NOT the bitboard with the occupancy of friendly pieces to avoid capturing yourself
var KNIGHT_MOVES [64]BitBoard

func initKnightMoves() {
	// UP-RIGHT +17
	// UP-LEFT +15
	// RIGHT-UP +10
	// RIGHT-DOWN -6
	// DOWN-RIGHT -15
	// DOWN-LEFT -17
	// LEFT-DOWN -10
	// LEFT-UP +6
	directions := []int{17, 15, 10, -6, -15, -17, -10, 6}
	for i := range NUM_SQUARES {
		for _, d := range directions {
			t := i + d // Target square

			// Check out of bounds (top/bottom)
			if t < 0 || t > 63 {
				continue
			}

			// Check out of bounds (crossing left/right)
			// Make sure the column shift was only 2 or less columns
			colShift := (t % 8) - (i % 8)
			if colShift < -2 || colShift > 2 {
				continue
			}

			KNIGHT_MOVES[i] |= Square(t).bitBoardPosition()
		}
	}
}
