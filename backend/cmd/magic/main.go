package main

import (
	"backend/internal/engine"

	"fmt"
	"math/bits"
	"math/rand/v2"
)

/*
This is to find magic numbers for rooks and bishops.
Just ran once and then copied the result into the code, but here for historic purposes
*/
func main() {
	fmt.Println("var ROOK_MAGICS = [64]uint64{")

	for sq := range engine.NUM_SQUARES {
		mask := engine.GenerateRookMask(engine.Square(sq))
		shift := 64 - bits.OnesCount64(uint64(mask))

		magic := findMagicNumber(sq, mask, shift, true) // true for Rook

		fmt.Printf("  0x%X, // Sq: %d, Bits: %d\n", magic, sq, 64-shift)
	}

	fmt.Println("}")

	fmt.Println("var BISHOP_MAGICS = [64]uint64{")

	for sq := range engine.NUM_SQUARES {
		mask := engine.GenerateBishopMask(engine.Square(sq))
		shift := 64 - bits.OnesCount64(uint64(mask))

		magic := findMagicNumber(sq, mask, shift, false) // false for Bishop

		fmt.Printf("  0x%X, // Sq: %d, Bits: %d\n", magic, sq, 64-shift)
	}

	fmt.Println("}")
}

// Sparse random numbers (few bits set) are statistically much more likely to be Magic Numbers
func randomMagicCandidate() uint64 {
	return rand.Uint64() & rand.Uint64() & rand.Uint64()
}

// Function to find magic numbers
func findMagicNumber(sq int, mask engine.BitBoard, shift int, isRook bool) uint64 {

	numBits := bits.OnesCount64(uint64(mask))
	variations := 1 << numBits

	occupancies := make([]engine.BitBoard, variations)
	correctMoves := make([]engine.BitBoard, variations)

	for i := range variations {
		occupancies[i] = engine.SetMaskOccupancy(i, numBits, mask)

		if isRook {
			correctMoves[i] = engine.GenerateSlowRookMoves(engine.Square(sq), occupancies[i])
		} else {
			correctMoves[i] = engine.GenerateSlowBishopMoves(engine.Square(sq), occupancies[i])
		}
	}

	// Try random numbers until one works.
	attempt := 0
	for {
		attempt++
		candidate := randomMagicCandidate()

		// This array simulates the magic hash table
		tableSize := 1 << (64 - shift)
		used := make([]engine.BitBoard, tableSize)
		filled := make([]bool, tableSize)

		// Test this candidate against ALL variations
		failed := false
		for i := range variations {

			// The magic formula
			idx := (uint64(occupancies[i]) * candidate) >> shift

			if !filled[idx] {
				used[idx] = correctMoves[i]
				filled[idx] = true
			} else if used[idx] != correctMoves[i] {
				// Failed to find a magic number
				failed = true
				break
			}
		}

		if !failed {
			// Found a magic number
			return candidate
		}
	}
}
