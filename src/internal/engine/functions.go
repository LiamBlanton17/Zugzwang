package engine

import (
	"fmt"
	"math/bits"
	"math/rand"
	"regexp"
	"strings"
)

/*
This file contains random helper functions that don't belong anywhere else
*/

// Converts a string (such as h7 or d3) to a Square (a number that will map to 0-63 on a bitboard)
func stringToSquare(s string) (Square, error) {
	// 255 is chosen to be the "NULL" value for the square
	if s == "-" || s == "" || s == " " {
		return Square(255), nil
	}

	// Once the "NULL" value is handled, square strings must have a length of 2
	// First is the column (a, b, c, etc) second is the row (1, 2, 3, etc)
	if len(s) != 2 {
		return Square(255), fmt.Errorf("Invalid square length of %v (should be 2 [col:row])", len(s))
	}

	// Check the column
	col := 0
	switch s[0] {
	case 'a':
		col = 0
	case 'b':
		col = 1
	case 'c':
		col = 2
	case 'd':
		col = 3
	case 'e':
		col = 4
	case 'f':
		col = 5
	case 'g':
		col = 6
	case 'h':
		col = 7
	default:
		return Square(255), fmt.Errorf("Invalid square column: %v", s[0])
	}

	// Check the row
	// This is done verbosely instead of converting string to int to check for more errors
	row := 0
	switch s[1] {
	case '1':
		row = 0
	case '2':
		row = 1
	case '3':
		row = 2
	case '4':
		row = 3
	case '5':
		row = 4
	case '6':
		row = 5
	case '7':
		row = 6
	case '8':
		row = 7
	default:
		return Square(255), fmt.Errorf("Invalid square row: %v", s[1])
	}

	// Combine row and column into one square
	return Square(row*8 + col), nil
}

// Converts a Square (a number that will map to 0-63 on a bitboard) to a string (such as h7 or d3)
func (sq Square) toString() string {
	// 255 is chosen to be the "NULL" value for the square
	if sq == NO_SQUARE {
		return " "
	}

	// Check the column
	str := ""
	switch sq % 8 {
	case 0:
		str = "a"
	case 1:
		str = "b"
	case 2:
		str = "c"
	case 3:
		str = "d"
	case 4:
		str = "e"
	case 5:
		str = "f"
	case 6:
		str = "g"
	case 7:
		str = "h"
	}

	// Check the row
	switch sq / 8 {
	case 0:
		str += "1"
	case 1:
		str += "2"
	case 2:
		str += "3"
	case 3:
		str += "4"
	case 4:
		str += "5"
	case 5:
		str += "6"
	case 6:
		str += "7"
	case 7:
		str += "8"
	}

	return str
}

// Converts a Piece to a character
func (p Piece) toString(c Color) string {
	switch p {
	case KING:
		if c == BLACK {
			return "k"
		}
		return "K"
	case QUEEN:
		if c == BLACK {
			return "q"
		}
		return "Q"
	case ROOK:
		if c == BLACK {
			return "r"
		}
		return "R"
	case BISHOP:
		if c == BLACK {
			return "b"
		}
		return "B"
	case KNIGHT:
		if c == BLACK {
			return "n"
		}
		return "N"
	case PAWN:
		if c == BLACK {
			return "p"
		}
		return "P"
	}

	return "-"
}

// Converts a Move to a string - allways sets the promotion piece to WHITE
func (m Move) toString() string {
	start := m.start.toString()
	target := m.target.toString()
	promotion := m.promotion.toString(WHITE)
	codeStr := "MOVE_CODE_NONE"
	switch m.code {
	case MOVE_CODE_CAPTURE:
		codeStr = "MOVE_CODE_CAPTURE"
	case MOVE_CODE_EN_PASSANT:
		codeStr = "MOVE_CODE_EN_PASSANT"
	case MOVE_CODE_DOUBLE_PAWN_PUSH:
		codeStr = "MOVE_CODE_DOUBLE_PAWN_PUSH"
	case MOVE_CODE_CASTLE:
		codeStr = "MOVE_CODE_CASTLE"
	}

	return fmt.Sprintf("(%v to %v) (promotion: %v) (code: %v)", start, target, promotion, codeStr)
}

// Converts a move to the pure cordniates notations (PCN)
func (m Move) toPCN() string {
	promo := ""
	if m.promotion != NO_PIECE {
		promo = m.promotion.toString(BLACK)
	}
	return fmt.Sprintf("%v%v%v", m.start.toString(), m.target.toString(), promo)
}

// Converts a string of SAN to PCN
// Regex to break down SAN:
// 1: Piece (NBKRQ or empty for pawn)
// 2: Disambiguation (file a-h or rank 1-8, or both)
// 3: Capture 'x' (optional, ignored)
// 4: Target Square (e.g. e4)
// 5: Promotion (e.g. =Q)
// 6: Check/Mate (+/#, optional)
var sanRegex = regexp.MustCompile(`^([NBKRQ])?([a-h1-8]{1,2})?x?([a-h][1-8])(=[NBRQ])?(\+|#)?$`)

func (b *Board) SanToPCN(san string) (string, error) {
	// 1. Handle Castling explicitly
	if san == "O-O" {
		return "e1g1", nil
	}
	if san == "O-O-O" {
		return "e1c1", nil
	}
	if san == "o-o" {
		return "e8g8", nil
	}
	if san == "o-o-o" {
		return "e8c8", nil
	}

	// 2. Parse the SAN string
	matches := sanRegex.FindStringSubmatch(san)
	if matches == nil {
		return "", fmt.Errorf("invalid SAN format: %s", san)
	}

	pieceChar := matches[1]      // "N", "B", etc. (Empty for Pawn)
	disambiguation := matches[2] // "a", "1", "bd", etc.
	targetStr := matches[3]      // "e4"
	promotionStr := matches[4]   // "=Q"

	// Convert target string to square index (0-63)
	targetSq, _ := stringToSquare(targetStr)

	// Determine the piece type we are looking for
	targetPieceType := PAWN
	if pieceChar != "" {
		switch pieceChar {
		case "N":
			targetPieceType = KNIGHT
		case "B":
			targetPieceType = BISHOP
		case "R":
			targetPieceType = ROOK
		case "Q":
			targetPieceType = QUEEN
		case "K":
			targetPieceType = KING
		}
	}

	// 3. Generate all legal moves to find the candidate
	// (Assuming you have a GenerateLegalMoves method)
	moves := make([]Move, 256)
	numberOfMoves := b.generatePseudoLegalMoves(moves)

	var candidate Move
	found := false

	for _, m := range moves[:numberOfMoves] {
		unmake, isLegal := b.makeMove(m)
		if !isLegal {
			b.unMakeMove(unmake)
			continue
		}
		b.unMakeMove(unmake)

		// A. Check Target Square
		if m.target != targetSq {
			continue
		}

		// B. Check Piece Type
		// We need to look at what piece is currently on the start square
		movingPiece := b.MailBox[m.start]
		// Extract type (0-5) from your internal Piece byte
		if movingPiece != targetPieceType {
			continue
		}

		// C. Check Promotion (if applicable)
		if promotionStr != "" {
			// promotionStr is "=Q", we need to check if move promotes to Queen
			wantedPromo := QUEEN
			switch promotionStr {
			case "=N":
				wantedPromo = KNIGHT
			case "=B":
				wantedPromo = BISHOP
			case "=R":
				wantedPromo = ROOK
			}
			if m.promotion != wantedPromo {
				continue
			}
		} else if m.promotion != NO_PIECE {
			// If SAN didn't specify promotion, but this move is a promotion, skip it
			// (e.g., SAN "e8" shouldn't match a move that is "e7e8q")
			continue
		}

		// D. Check Disambiguation (e.g. "Nbd7" -> start file must be 'b')
		if disambiguation != "" {
			startSqStr := m.start.toString()
			startFile := string(startSqStr[0])
			startRank := string(startSqStr[1])

			// If disambiguation contains file (letter)
			if strings.ContainsAny(disambiguation, "abcdefgh") {
				if !strings.Contains(startFile, disambiguation) && !strings.Contains(disambiguation, startFile) {
					// Checking containment handles full "e2" disambiguation cases
					if !strings.HasSuffix(disambiguation, startFile) && !strings.HasPrefix(disambiguation, startFile) {
						// Simple check: does the start file match the disambiguation char?
						if !strings.Contains(disambiguation, startFile) {
							continue
						}
					}
				}
			}
			// If disambiguation contains rank (number)
			if strings.ContainsAny(disambiguation, "12345678") {
				if !strings.Contains(disambiguation, startRank) {
					continue
				}
			}
		}

		// Found a match!
		candidate = m
		found = true
		break
	}

	if !found {
		return "", fmt.Errorf("no legal move found for %s", san)
	}

	// 4. Convert the found move to PCN string (e.g., "e2e4" or "a7a8q")
	startStr := candidate.start.toString()
	finalTargetStr := candidate.target.toString()
	promoSuffix := ""
	if candidate.promotion != NO_PIECE {
		// PCN usually uses lowercase for promotion (e.g., e7e8q)
		promoSuffix = strings.ToLower(candidate.promotion.toString(BLACK))
	}

	return startStr + finalTargetStr + promoSuffix, nil
}

// Get the move ordering score of the Move -- for move ordering
func (m *Move) orderScore(board *Board, ttEntry *TTEntry, killers *[2]Move, twoPlyKillers *[2]Move, cutoffHistory *CutoffHeuristic) int {

	// Check the TT table
	if ttEntry != nil && ttEntry.move == *m {
		return 1_000_000
	}

	// Check promotion
	if m.promotion != NO_PIECE {
		return 900_000 + int(PIECE_VALUES[m.promotion])
	}

	// Check captures MVV-LVA
	if m.code == MOVE_CODE_CAPTURE {
		return 800_000 + int((PIECE_VALUES[board.MailBox[m.target]]*10)-PIECE_VALUES[board.MailBox[m.start]])
	}

	// En passent is also a caputre
	if m.code == MOVE_CODE_EN_PASSANT {
		return 800_900
	}

	// Check killers
	// Place killers in right below equal MVV-LVA (PAWN takes PAWN is 800_900, BISHOP takes PAWN is 800_700)
	if killers != nil {
		if *m == (*killers)[0] {
			return 800_705
		}
		if *m == (*killers)[1] {
			return 800_704
		}
	}

	// Check the killer moves at a previous ply
	// Place killers in right below equal MVV-LVA
	if twoPlyKillers != nil {
		if *m == (*twoPlyKillers)[0] {
			return 800_703
		}
		if *m == (*twoPlyKillers)[1] {
			return 800_702
		}
	}

	// Check cutoff history
	// Cap history to prevent it from overtaking killers/captures
	score := 0
	if cutoffHistory != nil {
		score = min(cutoffHistory[board.Turn][m.start][m.target], 650_000)
	}

	// Castling bonus - boost castling above regular quiet moves
	if m.code == MOVE_CODE_CASTLE {
		return score + 7_000
	}

	return score
}

// Helper function to get the index on a bitboard given some square number
func (s Square) bitBoardPosition() BitBoard {
	return BitBoard(uint64(1) << s)
}

// Helper function to get the LSB (Square) of a bitboard and pop it off
func (bb *BitBoard) popSquare() Square {
	// Index of the rightmost set bit (Least Significant Bit)
	idx := bits.TrailingZeros64(uint64(*bb))

	// Removes the lowest set bit
	*bb &= (*bb - 1)

	return Square(idx)
}

// Helper function to clear a bit off a bitboard
func (bb *BitBoard) clear(b BitBoard) {
	*bb &= ^b
}

// Helper function to set a bit on a bitboard
func (bb *BitBoard) set(b BitBoard) {
	*bb |= b
}

// Helper function to setup magic bitboards
func SetMaskOccupancy(index int, bitsInMask int, attackMask BitBoard) BitBoard {
	occupancy := BitBoard(0)

	for i := range bitsInMask {
		square := attackMask.popSquare()

		// If the i-th bit of 'index' is set, place a piece there
		if (index & (1 << i)) != 0 {
			occupancy |= (1 << square)
		}
	}

	return occupancy
}

// Globals and function to setup for Zobrist hashing
const MASTER_ZOBRIST = 20240928                                    // Used for initializing Zobrist values
var PIECE_ZOBRIST [NUM_COLORS][NUM_PIECES][NUM_SQUARES]ZobristHash // Global for Zobrist hashing pieces
var BLACK_TO_MOVE_ZOBRIST ZobristHash                              // Global for Zorbist hasing black to move
var CASTLING_ZOBRIST [16]ZobristHash                               // Global for Zobrist hashing castling rights (1 for each combination)
var ENPASSENT_ZOBRIST [8]ZobristHash                               // Global for Zobrist hashing enpassent column (8 columns totoal)
func initZobrist() {
	// Setup determistic hashing with one constant master key
	source := rand.NewSource(MASTER_ZOBRIST)
	rng := rand.New(source)

	// Setup piece hashing
	for color := range NUM_COLORS {
		for piece := range NUM_PIECES {
			for square := range NUM_SQUARES {
				PIECE_ZOBRIST[color][piece][square] = ZobristHash(rng.Uint64())
			}
		}
	}

	// Setup to move hashing
	BLACK_TO_MOVE_ZOBRIST = ZobristHash(rng.Uint64())

	// Setup castling hashing
	for i := range CASTLING_ZOBRIST {
		CASTLING_ZOBRIST[i] = ZobristHash(rng.Uint64())
	}

	// Setup en passent hashing
	for i := range ENPASSENT_ZOBRIST {
		ENPASSENT_ZOBRIST[i] = ZobristHash(rng.Uint64())
	}
}
