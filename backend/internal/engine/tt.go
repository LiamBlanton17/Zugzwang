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
	zobrist ZobristHash
	eval    Eval
	depth   uint8
	flag    uint8
	unused  uint32 // padding out to 128 bits, will be used later
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
// With 512MB of memory, the size of the TT table can be ~33,554,432 entries if the tt entry size is 128 btis
const TT_SIZE = 33554432

// Global TT table
// This is to be used by all threads and games
var TT [TT_SIZE]TTEntry

// Function to probe the tt
func probeTT(zobrist ZobristHash) (TTEntry, bool) {
	// Compute the TT key
	key := zobrist & (TT_SIZE - 1)

	// Get the entry
	tt := TT[key]

	// See if the entry is there and valid
	if tt.flag != TT_LOCKED && tt.zobrist == zobrist {

		// Return a tt hit
		return tt, true
	}

	// Return a tt miss
	return tt, false
}

// Function to update the tt
func updateTT(zobrist ZobristHash, eval Eval, flag, depth uint8) {
	// Compute the TT key
	key := zobrist & (TT_SIZE - 1)

	// Lock the TT entry
	TT[key].flag = TT_LOCKED

	// Update the TT entry
	TT[key].zobrist = zobrist
	TT[key].eval = eval
	TT[key].depth = depth

	// Set the flag and unlock the entry
	TT[key].flag = flag
}
