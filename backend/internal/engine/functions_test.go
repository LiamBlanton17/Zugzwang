package engine

import "testing"

// Approach will be to valide "NULL" squares, and then go down a diagonal making sure each column and row parses correctly
func TestStringToSquare(t *testing.T) {
	// Tests setup to be run
	tests := []struct {
		name        string
		input       string
		result      Square
		expectError bool
	}{
		{
			name:        "Validate empty #1",
			input:       "",
			result:      255,
			expectError: false,
		},
		{
			name:        "Validate empty #2",
			input:       "-",
			result:      255,
			expectError: false,
		},
		{
			name:        "Validate empty #3",
			input:       " ",
			result:      255,
			expectError: false,
		},
		{
			name:        "Verify a1 parses correctly",
			input:       "a1",
			result:      0,
			expectError: false,
		},
		{
			name:        "Verify b2 parses correctly",
			input:       "b2",
			result:      9,
			expectError: false,
		},
		{
			name:        "Verify c3 parses correctly",
			input:       "c3",
			result:      18,
			expectError: false,
		},
		{
			name:        "Verify d4 parses correctly",
			input:       "d4",
			result:      27,
			expectError: false,
		},
		{
			name:        "Verify e5 parses correctly",
			input:       "e5",
			result:      36,
			expectError: false,
		},
		{
			name:        "Verify f6 parses correctly",
			input:       "f6",
			result:      45,
			expectError: false,
		},
		{
			name:        "Verify g7 parses correctly",
			input:       "g7",
			result:      54,
			expectError: false,
		},
		{
			name:        "Verify h8 parses correctly",
			input:       "h8",
			result:      63,
			expectError: false,
		},
		{
			name:        "Verify u9 fails",
			input:       "u9",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify yyy fails",
			input:       "yyy",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify 12345 fails",
			input:       "12345",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify a fails",
			input:       "a",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify 4 fails",
			input:       "4",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify A1 fails",
			input:       "A1",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify H8 fails",
			input:       "H8",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify a0 fails",
			input:       "a0",
			result:      255,
			expectError: true,
		},
		{
			name:        "Verify h9 fails",
			input:       "h9",
			result:      255,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := stringToSquare(tc.input)

			// Check if error was expected but did not happen
			if tc.expectError && err == nil {
				t.Errorf("Expected an error for sqaure %q, but got none", tc.input)
			}

			// Check if error was not expected but did happen
			if !tc.expectError && err != nil {
				t.Errorf("Did not expect error for sqaure %q, but got: %v", tc.input, err)
			}

			// Verify the results match (if error was not expected)
			if !tc.expectError && tc.result != result {
				t.Errorf("Expected %v but incorrectly got %v", tc.result, result)
			}
		})
	}
}

func TestBitBoardPosition(t *testing.T) {
	// Tests setup to be run
	tests := []struct {
		name   string
		input  Square
		result BitBoard
	}{
		{
			name:   "Test #1",
			input:  63,
			result: 0b1000000000000000000000000000000000000000000000000000000000000000,
		},
		{
			name:   "Test #2",
			input:  0,
			result: 0b0000000000000000000000000000000000000000000000000000000000000001,
		},
		{
			name:   "Test #3",
			input:  17,
			result: 0b0000000000000000000000000000000000000000000000100000000000000000,
		},
		{
			name:   "Test #4",
			input:  36,
			result: 0b0000000000000000000000000001000000000000000000000000000000000000,
		},
		{
			name:   "Test #5",
			input:  55,
			result: 0b0000000010000000000000000000000000000000000000000000000000000000,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.bitBoardPosition()

			// Verify the results match
			if tc.result != result {
				t.Errorf("Expected %v but incorrectly got %v", tc.result, result)
			}
		})
	}
}
