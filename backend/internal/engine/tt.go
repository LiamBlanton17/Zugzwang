package engine

/*
This file holds the logic to the transposition tables for the chess engine
On the cloud server, memory will be highly constrained so they should be easily edittable
*/

/*
TODO: Redo the locking on the TT tables to use atomic reads/writes and be non-locking, once on production
*/

// TT entry definition (currently 128 bits)
type TTEntry struct {
	move    Move
	eval    Eval
	depth   uint8
	flag    uint8
	zobrist ZobristHash
}

// TT flags
const (
	TT_LOCKED uint8 = iota
	TT_EXACT
	TT_UPPER
	TT_LOWER
)

// Size of the TT table
// The memory will be TT_SIZE * sizeof(TTEntry)
// With 256MB of memory, the size of the TT table can be ~16,777,216 entries if the tt entry size is 128 btis
const TT_SIZE = 16777216

// Use a slice instead of a fixed-size array
var TT []TTEntry

func initTT() {
	TT = make([]TTEntry, TT_SIZE)
}

// Now this will work:
func ClearTT() {
	clear(TT)
}

// Function to probe the tt
func probeTT(zobrist ZobristHash) (*TTEntry, bool) {
	// Compute the TT key
	key := zobrist & (TT_SIZE - 1)

	// Get the entry and return hit
	entry := &TT[key]
	if entry.zobrist == zobrist {
		return entry, true
	}

	// Return miss
	return entry, false
}

// Function to update the tt
func updateTT(zobrist ZobristHash, eval Eval, flag, depth uint8, move Move) {
	// Compute the TT key
	key := zobrist & (TT_SIZE - 1)

	// Get the TT entry
	entry := &TT[key]

	// Write everything except zobrist
	entry.eval = eval
	entry.depth = depth
	entry.flag = flag
	entry.move = move

	// Write the Zobrist update last
	entry.zobrist = zobrist
}
