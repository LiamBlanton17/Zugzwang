package engine

import (
	"fmt"
	"math/bits"
)

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
// This is inefficient
// Should be more like pushes := (whitePawns << 8) & ^occupancy to get all the pawn moves for white one push
// Todo: Refactor later to make more efficient
func getPawnMoves(pawns, enemyPieces, occupancy BitBoard, EPS Square, color Color, moves []Move, moveIdx int) int {

	for pawns > 0 {
		start := pawns.popSquare()

		if color == WHITE { // WHITE pawn moves
			oneSq := start + 8
			twoSq := start + 16

			// Single Push (Not promotions, those are handled at the end of the code)
			if occupancy&oneSq.bitBoardPosition() == 0 {
				if oneSq <= 55 {
					moveIdx = addMove(moves, start, oneSq, MOVE_CODE_NONE, false, moveIdx)

					// Double Push (on second rank)
					if start >= 8 && start <= 15 && (occupancy&twoSq.bitBoardPosition()) == 0 {
						moveIdx = addMove(moves, start, twoSq, MOVE_CODE_DOUBLE_PAWN_PUSH, false, moveIdx)
					}
				}
			}

			// Capture targets
			capLeft := start + 7
			capRight := start + 9

			// Prevent wrapping around board by making sure start column is not the A column (col 0)
			canCapLeft := start%8 > 0 && ((enemyPieces&capLeft.bitBoardPosition()) != 0 || capLeft == EPS)

			// Prevent wrapping around board by making sure start column is not the H column (col 7)
			canCapRight := start%8 < 7 && ((enemyPieces&capRight.bitBoardPosition()) != 0 || capRight == EPS)

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
				if (occupancy & oneSq.bitBoardPosition()) == 0 {
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
			if (occupancy & oneSq.bitBoardPosition()) == 0 {
				if oneSq >= 8 {
					moveIdx = addMove(moves, start, oneSq, MOVE_CODE_NONE, false, moveIdx)
					if start >= 48 && start <= 55 && (occupancy&twoSq.bitBoardPosition()) == 0 {
						moveIdx = addMove(moves, start, twoSq, MOVE_CODE_DOUBLE_PAWN_PUSH, false, moveIdx)
					}
				}
			}

			// Capture targets
			capRight := start - 7
			capLeft := start - 9

			// Prevent wrapping around board by making sure start column is not the A column (col 0)
			canCapRight := start%8 < 7 && ((enemyPieces&capRight.bitBoardPosition()) != 0 || capRight == EPS)

			// Prevent wrapping around board by making sure start column is not the H column (col 7)
			canCapLeft := start%8 > 0 && ((enemyPieces&capLeft.bitBoardPosition()) != 0 || capLeft == EPS)

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
				if occupancy&oneSq.bitBoardPosition() == 0 {
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

// Global magic lookup table for Rooks (reused for queens)
// 4096 comes from the number 2^12, where 12 is the number of bits we mask per rook position
// During move generation at run time, AND NOT the bitboard with the occupancy of friendly pieces to avoid capturing yourself
var MAGIC_ROOK_MOVES [NUM_SQUARES][4096]BitBoard

// Used for actually looking up moves in these magic tables
type Magic struct {
	mask  BitBoard
	magic uint64
	shift int
}

func (m *Magic) getMagicIndex(occupancy BitBoard) int {
	return int((uint64(occupancy&m.mask) * m.magic) >> m.shift)
}

// Global rook magic info
var MAGIC_ROOK_INFO [NUM_SQUARES]Magic

func initMagicRook() {

	for i := range NUM_SQUARES {

		// Initialize the masking
		mask := GenerateRookMask(Square(i))
		MAGIC_ROOK_INFO[i].mask = mask

		// Setup the magic number from magics.go
		MAGIC_ROOK_INFO[i].magic = ROOK_MAGICS[i]

		// Given the mask, initialize the magic board
		// Count bits in the mask
		numBits := bits.OnesCount64(uint64(mask))
		MAGIC_ROOK_INFO[i].shift = 64 - numBits

		// Calculate total variations
		variations := 1 << numBits

		// Loop through every possible occupancy pattern
		for j := range variations {
			// Map the index 'j' to an actual board of pieces
			occupancy := SetMaskOccupancy(j, numBits, mask)

			// Calculate the magic index and get the moves
			magicIndex := MAGIC_ROOK_INFO[i].getMagicIndex(occupancy)
			moves := GenerateSlowRookMoves(Square(i), occupancy)

			// Need to panic if Magic numbers fail to init the magic tables correctly
			// Maybe this can be handled more gracefully, but is mostly fine because magics should be
			// hard coded correctly, and if they are not need to be fixed
			if MAGIC_ROOK_MOVES[i][magicIndex] != 0 && MAGIC_ROOK_MOVES[i][magicIndex] != moves {
				panic(fmt.Sprintf("Magic Collision on Square %d! Index %d used twice for different moves.", i, magicIndex))
			}

			// Store the moves
			MAGIC_ROOK_MOVES[i][magicIndex] = moves
		}
	}

}

// Get the mask of relevant bits from a rook square
func GenerateRookMask(sq Square) BitBoard {
	// Initialize the masking
	mask := BitBoard(0)

	// Go up
	// Ignore top row
	j := Square(sq + 8)
	for j <= 55 {
		mask |= j.bitBoardPosition()
		j += 8
	}

	// Go down
	// Ignore last row
	j = Square(sq - 8)
	for j >= 8 {
		mask |= j.bitBoardPosition()
		j -= 8
	}

	// Go right (checking columns to prevent wrap arround)
	// Also, ignore the final column
	j = Square(sq + 1)
	if j%8 != 0 { // check that is hasn't wrapped around
		for j%8 < 7 {
			mask |= j.bitBoardPosition()
			j += 1
		}
	}

	// Go left (checking columns to prevent wrap arround)
	// Also, ignore the final column
	j = Square(sq - 1)
	if j%8 != 7 { // check that is hasn't wrapped around
		for j%8 > 0 {
			mask |= j.bitBoardPosition()
			j -= 1
		}
	}

	return mask
}

// This function should be called during the initialization of the magic rook bitboards, but not after
func GenerateSlowRookMoves(sq Square, occupancy BitBoard) BitBoard {
	moves := BitBoard(0)

	// Try to go up
	i := sq + 8
	for i < 64 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i += 8
	}

	// Try to go down
	// Super counterintuative, but i -= 8 will wrap around back to 250 something, so check under 64
	i = sq - 8
	for i < 64 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i -= 8
	}

	// Try to go right
	i = sq + 1
	for i%8 > 0 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i += 1
	}

	// Try to go left
	i = sq - 1
	for i%8 < 7 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i -= 1
	}

	return moves
}

// Global magic lookup table for Bishops (reused for queens)
// 512 comes from the number 2^9, where 9 is the number of bits we mask per bishop position
// During move generation at run time, AND NOT the bitboard with the occupancy of friendly pieces to avoid capturing yourself
var MAGIC_BISHOP_MOVES [NUM_SQUARES][512]BitBoard

// Used for occupancy masking
var MAGIC_BISHOP_MASK [NUM_SQUARES]BitBoard

// Global rook magic info
var MAGIC_BISHOP_INFO [NUM_SQUARES]Magic

func initMagicBishop() {

	for i := range NUM_SQUARES {

		// Initialize the masking
		mask := GenerateBishopMask(Square(i))
		MAGIC_BISHOP_INFO[i].mask = mask

		// Setup the magic number from magics.go
		MAGIC_BISHOP_INFO[i].magic = BISHOP_MAGICS[i]

		// Given the mask, initialize the magic board
		// Count bits in the mask
		numBits := bits.OnesCount64(uint64(mask))
		MAGIC_BISHOP_INFO[i].shift = 64 - numBits

		// Calculate total variations
		variations := 1 << numBits

		// Loop through every possible occupancy pattern
		for j := range variations {
			// Map the index 'j' to an actual board of pieces
			occupancy := SetMaskOccupancy(j, numBits, mask)

			// Calculate the magic index and get the moves
			magicIndex := MAGIC_BISHOP_INFO[i].getMagicIndex(occupancy)
			moves := GenerateSlowBishopMoves(Square(i), occupancy)

			// Need to panic if Magic numbers fail to init the magic tables correctly
			// Maybe this can be handled more gracefully, but is mostly fine because magics should be
			// hard coded correctly, and if they are not need to be fixed
			if MAGIC_BISHOP_MOVES[i][magicIndex] != 0 && MAGIC_BISHOP_MOVES[i][magicIndex] != moves {
				panic(fmt.Sprintf("Magic Collision on Square %d! Index %d used twice for different moves.", i, magicIndex))
			}

			// Store the moves
			MAGIC_BISHOP_MOVES[i][magicIndex] = moves
		}
	}

}

// Get the mask of relevant bits from a bishop square
func GenerateBishopMask(sq Square) BitBoard {
	// Initialize the masking
	mask := BitBoard(0)

	// Go up right
	// Ignore last column/row
	j := Square(sq + 9)
	if j%8 != 0 { // check that is hasn't wrapped around

		// Check that we are not final row and not final column
		for j <= 55 && j%8 < 7 {
			mask |= j.bitBoardPosition()
			j += 9
		}
	}

	// Go up left
	// Ignore last column/row
	j = Square(sq + 7)
	if j%8 != 7 { // check that is hasn't wrapped around

		// Check that we are not final row and not final column
		for j <= 55 && j%8 > 0 {
			mask |= j.bitBoardPosition()
			j += 7
		}
	}

	// Go down right
	// Ignore last column/row
	j = Square(sq - 7)
	if j%8 != 0 { // check that is hasn't wrapped around

		// Check that we are not final row and not final column
		for j >= 8 && j%8 < 7 {
			mask |= j.bitBoardPosition()
			j -= 7
		}
	}

	// Go down left
	// Ignore last column/row
	j = Square(sq - 9)
	if j%8 != 7 { // check that is hasn't wrapped around

		// Check that we are not final row and not final column
		for j >= 8 && j%8 > 0 {
			mask |= j.bitBoardPosition()
			j -= 9
		}
	}

	return mask
}

// This function should be called during the initialization of the magic bishop bitboards, but not after
func GenerateSlowBishopMoves(sq Square, occupancy BitBoard) BitBoard {
	moves := BitBoard(0)

	// Go up right
	i := sq + 9
	for i < 64 && i%8 > 0 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i += 9
	}

	// Go up left
	i = sq + 7
	for i < 64 && i%8 < 7 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i += 7
	}

	// Go down right
	i = sq - 7
	for i < 64 && i%8 > 0 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i -= 7
	}

	// Go down left
	i = sq - 9
	for i < 64 && i%8 < 7 {
		moves |= i.bitBoardPosition()
		if occupancy&i.bitBoardPosition() != 0 {
			break
		}

		i -= 9
	}

	return moves
}

// Global King lookup table
// Knight move generation logic is simple, can only move 1 direction each way
// Castling is handled elsewhere
// During move generation at run time, AND NOT the bitboard with the occupancy of friendly pieces to avoid capturing yourself
var KING_MOVES [NUM_SQUARES]BitBoard

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
var KNIGHT_MOVES [NUM_SQUARES]BitBoard

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
