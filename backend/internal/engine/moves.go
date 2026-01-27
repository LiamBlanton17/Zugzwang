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

// Used by board.generateMoves() to get the pseudo-legal queen moves
func (b *Board) getQueenMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	queens := b.Pieces[b.Turn][QUEEN]
	enemyPieces := b.getEnemyPieces()

	for queens > 0 {
		start := queens.popSquare()
		magicIdx := MAGIC_ROOK_INFO[start].getMagicIndex(b.Occupancy[EITHER_COLOR])
		magicMoves := MAGIC_ROOK_MOVES[start][magicIdx] &^ b.Occupancy[b.Turn]

		for magicMoves > 0 {
			target := magicMoves.popSquare()

			// Handle move code if a capture
			code := MOVE_CODE_NONE
			if target.bitBoardPosition()&enemyPieces != 0 {
				code = MOVE_CODE_CAPTURE
			}
			moveIdx = addMove(moves, start, target, code, false, moveIdx)
		}

		magicIdx = MAGIC_BISHOP_INFO[start].getMagicIndex(b.Occupancy[EITHER_COLOR])
		magicMoves = MAGIC_BISHOP_MOVES[start][magicIdx] &^ b.Occupancy[b.Turn]

		for magicMoves > 0 {
			target := magicMoves.popSquare()

			// Handle move code if a capture
			code := MOVE_CODE_NONE
			if target.bitBoardPosition()&enemyPieces != 0 {
				code = MOVE_CODE_CAPTURE
			}
			moveIdx = addMove(moves, start, target, code, false, moveIdx)
		}
	}

	b.MoveIdx = moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal rook moves
func (b *Board) getRookMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	rooks := b.Pieces[b.Turn][ROOK]
	enemyPieces := b.getEnemyPieces()

	for rooks > 0 {
		start := rooks.popSquare()
		magicIdx := MAGIC_ROOK_INFO[start].getMagicIndex(b.Occupancy[EITHER_COLOR])
		magicMoves := MAGIC_ROOK_MOVES[start][magicIdx] &^ b.Occupancy[b.Turn]

		for magicMoves > 0 {
			target := magicMoves.popSquare()

			// Handle move code if a capture
			code := MOVE_CODE_NONE
			if target.bitBoardPosition()&enemyPieces != 0 {
				code = MOVE_CODE_CAPTURE
			}
			moveIdx = addMove(moves, start, target, code, false, moveIdx)
		}
	}

	b.MoveIdx = moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal bishop moves
func (b *Board) getBishopMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	bishops := b.Pieces[b.Turn][BISHOP]
	enemyPieces := b.getEnemyPieces()

	for bishops > 0 {
		start := bishops.popSquare()
		magicIdx := MAGIC_BISHOP_INFO[start].getMagicIndex(b.Occupancy[EITHER_COLOR])
		magicMoves := MAGIC_BISHOP_MOVES[start][magicIdx] &^ b.Occupancy[b.Turn]

		for magicMoves > 0 {
			target := magicMoves.popSquare()

			// Handle move code if a capture
			code := MOVE_CODE_NONE
			if target.bitBoardPosition()&enemyPieces != 0 {
				code = MOVE_CODE_CAPTURE
			}
			moveIdx = addMove(moves, start, target, code, false, moveIdx)
		}
	}

	b.MoveIdx = moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal pawn moves
// Does include captures
// Does include promotion
// Does include enpassent
// This is inefficient
// Should be more like pushes := (whitePawns << 8) & ^occupancy to get all the pawn moves for white one push
// Todo: Refactor later to make more efficient
func (b *Board) getPawnMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	pawns := b.Pieces[b.Turn][PAWN]
	occupancy := b.Occupancy[b.Turn]
	enemyPieces := b.getEnemyPieces()

	for pawns > 0 {
		start := pawns.popSquare()

		if b.Turn == WHITE { // WHITE pawn moves
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
			canCapLeft := start%8 > 0 && ((enemyPieces&capLeft.bitBoardPosition()) != 0 || capLeft == b.EPS)

			// Prevent wrapping around board by making sure start column is not the H column (col 7)
			canCapRight := start%8 < 7 && ((enemyPieces&capRight.bitBoardPosition()) != 0 || capRight == b.EPS)

			// Captures (Not promotion)
			if oneSq <= 55 {
				if canCapLeft {
					code := MOVE_CODE_CAPTURE
					if capLeft == b.EPS {
						code = MOVE_CODE_EN_PASSANT
					}
					moveIdx = addMove(moves, start, capLeft, code, false, moveIdx)
				}
				if canCapRight {
					code := MOVE_CODE_CAPTURE
					if capRight == b.EPS {
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
			canCapRight := start%8 < 7 && ((enemyPieces&capRight.bitBoardPosition()) != 0 || capRight == b.EPS)

			// Prevent wrapping around board by making sure start column is not the H column (col 7)
			canCapLeft := start%8 > 0 && ((enemyPieces&capLeft.bitBoardPosition()) != 0 || capLeft == b.EPS)

			// Captures (Not promotion)
			if oneSq >= 8 {
				if canCapLeft {
					code := MOVE_CODE_CAPTURE
					if capLeft == b.EPS {
						code = MOVE_CODE_EN_PASSANT
					}
					moveIdx = addMove(moves, start, capLeft, code, false, moveIdx)
				}
				if canCapRight {
					code := MOVE_CODE_CAPTURE
					if capRight == b.EPS {
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

	b.MoveIdx = moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal king moves
// This does not include castling
func (b *Board) getKingMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	king := b.Pieces[b.Turn][KING]

	for king > 0 {
		start := king.popSquare()
		targets := KING_MOVES[start] &^ b.Occupancy[b.Turn]

		for targets > 0 {
			moveIdx = addMove(moves, start, targets.popSquare(), MOVE_CODE_NONE, false, moveIdx)
		}
	}

	b.MoveIdx = moveIdx
}

// Used by board.generateMoves() to get the pseudo-legal knight moves
func (b *Board) getKnightMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	knights := b.Pieces[b.Turn][KNIGHT]

	for knights > 0 {
		start := knights.popSquare()
		targets := KNIGHT_MOVES[start] &^ b.Occupancy[b.Turn]

		for targets > 0 {
			moveIdx = addMove(moves, start, targets.popSquare(), MOVE_CODE_NONE, false, moveIdx)
		}
	}

	b.MoveIdx = moveIdx
}

// Used by board.generateMoves() to get the legal castling moves
// This function checks for king is attacked after the move, which is the main logical difference between legal and pseudo-legal
func (b *Board) getCastlingMoves() {
	moves := b.Moves
	moveIdx := b.MoveIdx
	occupancy := b.Occupancy[EITHER_COLOR]

	// White can castle kingside
	if CASTLE_WK&b.CR != 0 && b.Turn == WHITE {

		// Make sure f1 and g1 are empty
		if F1.bitBoardPosition()&occupancy == 0 && G1.bitBoardPosition()&occupancy == 0 {

			// Make sure e1, f1 are not attacked (can't castle out of or through check)
			// Checking if the king is placed in check on g1 is handled later, when validating all moves against
			// Illegally putting the king in check
			if !b.isSquareAttacked(E1, BLACK) && !b.isSquareAttacked(F1, BLACK) {
				moveIdx = addMove(moves, E1, G1, MOVE_CODE_CASTLE, false, moveIdx)
			}
		}
	}

	// White can castle queenside
	if CASTLE_WQ&b.CR != 0 && b.Turn == WHITE {

		// Make sure d1, c1 and b1 are empty
		if D1.bitBoardPosition()&occupancy == 0 && C1.bitBoardPosition()&occupancy == 0 && B1.bitBoardPosition()&occupancy == 0 {

			// Make sure e1, d1 are not attacked (can't castle out of or through check)
			// Checking if the king is placed in check on c1 is handled later, when validating all moves against
			// Illegally putting the king in check
			// TODO: Optimize to avoid repeat calls checking e1 for both queen/kingside castling, though its probably not often
			if !b.isSquareAttacked(E1, BLACK) && !b.isSquareAttacked(D1, BLACK) {
				moveIdx = addMove(moves, E1, C1, MOVE_CODE_CASTLE, false, moveIdx)
			}
		}
	}

	// Black can castle kingside
	if CASTLE_BK&b.CR != 0 && b.Turn == BLACK {

		// Make sure f8 and g8 are empty
		if F8.bitBoardPosition()&occupancy == 0 && G8.bitBoardPosition()&occupancy == 0 {

			// Make sure e8, f8 are not attacked (can't castle out of or through check)
			// Checking if the king is placed in check on g8 is handled later, when validating all moves against
			// Illegally putting the king in check
			if !b.isSquareAttacked(E8, WHITE) && !b.isSquareAttacked(F8, WHITE) {
				moveIdx = addMove(moves, E8, G8, MOVE_CODE_CASTLE, false, moveIdx)
			}
		}
	}

	// Black can castle queenside
	if CASTLE_BQ&b.CR != 0 && b.Turn == BLACK {

		// Make sure b8, c8, and d8 are empty
		if D8.bitBoardPosition()&occupancy == 0 && C8.bitBoardPosition()&occupancy == 0 && B8.bitBoardPosition()&occupancy == 0 {

			// Make sure e8, d8  are not attacked (can't castle out of or through check)
			// Checking if the king is placed in check on c8 is handled later, when validating all moves against
			// Illegally putting the king in check
			// TODO: Optimize to avoid repeat calls checking E8 for both queen/kingside castling, though its probably not often
			if !b.isSquareAttacked(E8, WHITE) && !b.isSquareAttacked(D8, WHITE) {
				moveIdx = addMove(moves, E8, C8, MOVE_CODE_CASTLE, false, moveIdx)
			}
		}

	}

	b.MoveIdx = moveIdx
}

// Helper function to check if a square is under attack, most useful for checking if king is under attack after a pseudo-legal move
func (b *Board) isSquareAttacked(sq Square, attackerSide Color) bool {
	// Check Pawn Attacks
	if attackerSide == WHITE {
		if sq%8 > 0 { // Not on A column
			if sq >= 9 && (sq-9).bitBoardPosition()&b.Pieces[WHITE][PAWN] != 0 {
				return true
			}
		}
		if sq%8 < 7 { // Not on H column
			if sq >= 7 && (sq-7).bitBoardPosition()&b.Pieces[WHITE][PAWN] != 0 {
				return true
			}
		}
	} else {
		if sq%8 > 0 { // Not on A column
			if sq <= 56 && (sq+7).bitBoardPosition()&b.Pieces[BLACK][PAWN] != 0 {
				return true
			}
		}
		if sq%8 < 7 { // Not on H column
			if sq <= 54 && (sq+9).bitBoardPosition()&b.Pieces[BLACK][PAWN] != 0 {
				return true
			}
		}
	}

	// Check Knight Attacks
	if (KNIGHT_MOVES[sq] & b.Pieces[attackerSide][KNIGHT]) != 0 {
		return true
	}

	// Check King Attacks
	if (KING_MOVES[sq] & b.Pieces[attackerSide][KING]) != 0 {
		return true
	}

	// Check Bishop/Queen Diagonal Attacks
	// Reuse your magic bitboards!
	bIdx := MAGIC_BISHOP_INFO[sq].getMagicIndex(b.Occupancy[EITHER_COLOR])
	if (MAGIC_BISHOP_MOVES[sq][bIdx] & (b.Pieces[attackerSide][BISHOP] | b.Pieces[attackerSide][QUEEN])) != 0 {
		return true
	}

	// Check Rook/Queen Straight Attacks
	rIdx := MAGIC_ROOK_INFO[sq].getMagicIndex(b.Occupancy[EITHER_COLOR])
	if (MAGIC_ROOK_MOVES[sq][rIdx] & (b.Pieces[attackerSide][ROOK] | b.Pieces[attackerSide][QUEEN])) != 0 {
		return true
	}

	return false
}

// Helper function to add a move and increment the counter
func addMove(moves []Move, start, target Square, code uint8, is_promotion bool, moveIdx uint8) uint8 {
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

// Helper function to get the magic index
func (m *Magic) getMagicIndex(occupancy BitBoard) int {
	return int((uint64(occupancy&m.mask) * m.magic) >> m.shift)
}

// Global rook magic info
var MAGIC_ROOK_INFO [NUM_SQUARES]Magic

// Function should be called once at startup to initialize the magic rook tables
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

// Function should be called once at startup to initialize the magic bishop tables
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

// This function makes a move, in-place, on a board, and returns if that move was legal or not
func (b *Board) makeMove(move Move) bool {
	return true
}

// This function unmakes a move, in-place, on a board
func (b *Board) unMakeMove(move Move) {

}

// This function simply checks if a move was legal, utilizing make and unmake moves
func (b *Board) isMoveLegal(move Move) bool {
	isLegal := b.makeMove(move)
	b.unMakeMove(move)
	return isLegal
}
