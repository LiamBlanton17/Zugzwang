package engine

import (
	"testing"
)

func TestInitKingMoves(t *testing.T) {
	// Setup the global table
	initKingMoves()

	// Tests setup to be run
	tests := []struct {
		name        string
		squareIndex int
		expectedBB  uint64
	}{
		{
			name:        "Corner a1 (Index 0)",
			squareIndex: 0,
			// Expect moves to b1(1), b2(9), a2(8)
			expectedBB: (1 << 1) | (1 << 9) | (1 << 8),
		},
		{
			name:        "Corner h8 (Index 63)",
			squareIndex: 63,
			// Expect moves to g8(62), g7(54), h7(55)
			expectedBB: (1 << 62) | (1 << 54) | (1 << 55),
		},
		{
			name:        "Edge a4 (Index 24)",
			squareIndex: 24, // Left edge, middle of board
			// Expect 5 moves: a5(32), b5(33), b4(25), b3(17), a3(16)
			expectedBB: (1 << 32) | (1 << 33) | (1 << 25) | (1 << 17) | (1 << 16),
		},
		{
			name:        "Center d4 (Index 27)",
			squareIndex: 27,
			// Expect all 8 surrounding squares:
			// Top: c5(34), d5(35), e5(36)
			// Mid: c4(26),       , e4(28)
			// Bot: c3(18), d3(19), e3(20)
			expectedBB: (1 << 34) | (1 << 35) | (1 << 36) | (1 << 26) | (1 << 28) | (1 << 18) | (1 << 19) | (1 << 20),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Retrieve the generated bitboard
			result := uint64(KING_MOVES[tc.squareIndex])

			if result != tc.expectedBB {
				t.Errorf("Square %d (%s):\nExpected: %064b\nGot:      %064b",
					tc.squareIndex, tc.name, tc.expectedBB, result)
			}
		})
	}
}

func TestInitKnightMoves(t *testing.T) {
	// Setup the global table
	initKnightMoves()

	// Tests setup to be run
	tests := []struct {
		name        string
		squareIndex int
		expectedBB  uint64
	}{
		{
			name:        "Corner a1 (Index 0)",
			squareIndex: 0,
			// Expect moves to b3 (17) and c2 (10)
			expectedBB: (1 << 17) | (1 << 10),
		},
		{
			name:        "Corner h8 (Index 63)",
			squareIndex: 63,
			// Expect moves to f7 (53) and g6 (46)
			expectedBB: (1 << 53) | (1 << 46),
		},
		{
			name:        "Edge h1 (Index 7)",
			squareIndex: 7,
			expectedBB:  (1 << 13) | (1 << 22),
		},
		{
			name:        "Center d4 (Index 27)",
			squareIndex: 27,
			// Expect moves to e6 (44), c6 (42), f5 (37), b5 (33), b3(17), f3(21), c2(10), e2(12)
			expectedBB: (1 << 44) | (1 << 42) | (1 << 37) | (1 << 33) | (1 << 17) | (1 << 21) | (1 << 10) | (1 << 12),
		},
		{
			name:        "Near Edge b1 (Index 1)",
			squareIndex: 1,
			// Moves: a3 (16), c3 (18), d2 (11)
			expectedBB: (1 << 16) | (1 << 18) | (1 << 11),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resultBitboard := KNIGHT_MOVES[tc.squareIndex]

			// Verify the exact BitBoard
			if tc.expectedBB != 0 && uint64(resultBitboard) != tc.expectedBB {
				t.Errorf("Square %d: \nExpected Bitboard: %064b (Hex: 0x%x)\nActual Bitboard: %064b (Hex: 0x%x)", tc.squareIndex, tc.expectedBB, tc.expectedBB, uint64(resultBitboard), uint64(resultBitboard))
			}
		})
	}
}
