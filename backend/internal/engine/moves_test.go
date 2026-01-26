package engine

import (
	"testing"
)

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
