# Zugzwang

A chess engine written in Go from scratch with a focus on performance, correctness, and low-level bit manipulation. Estimated playing strength around 1600-1800 (in the top 98%+ of humans).

## Key Features

1. **Negamax with alpha-beta pruning and move ordering heuristics**

    The Zugzwang engine uses the negamax algorithm for searching, which is a variant of minimax. With alpha-beta pruning and several move ordering heuristics (such as killer moves and a transposition table), the search reaches a depth of 8 ply, plus an additional quiescence search, in about 1-2 seconds. Given that the branching factor of chess is around 30, a naive search to that depth would explore a state space of at least 656 billion positions — likely 5-10x more in practice. At Zugzwang's raw speed of approximately 1M nodes per second, we can estimate its effective branching factor to be around 7-8.

2. **Bitboard representation**

    Zugzwang uses 12 bitboards, one for each color and piece-type combination. This allows for extremely fast bitwise operations, leading to fast move generation and board state updates. Zugzwang also trades a bit of space for time by redundantly maintaining a "mailbox" array mapping squares to pieces, which enables efficient square-centric lookups during evaluation.

3. **Magic bitboards**

    Zugzwang uses magic bitboards to efficiently generate sliding-piece moves. Rather than iterating ray-by-ray for each sliding piece every time moves are generated, Zugzwang indexes into a precalculated table of attack sets. This leads to much faster move generation, resulting in greater search depth and more playing strength.

4. **Zobrist hashing**

    The engine uses Zobrist hashing to represent the entire board state as a single 64-bit integer. Updating the Zobrist hash after a move is a handful of XOR operations rather than a full recalculation. This hash keys a transposition table of previously evaluated positions, which prunes the search tree by avoiding redundant work on transpositions.

5. **Zero-allocation searching**

    The engine performs no heap allocations during the search. All memory used by the search is pre-allocated before it begins, keeping the hot path allocation-free.

6. **Custom evaluation function**

    Zugzwang includes a custom evaluation function designed around my own understanding of chess (I am around 1800-1900 ELO). It includes piece-square tables that are tapered between opening and endgame scores based on the material remaining on the board, file-based evaluation, king safety, and pawn structure terms.

7. **[WIP] Simple, stateless HTTP API**

    Zugzwang will expose a simple, stateless API with a single POST endpoint, intended for wiring up a frontend that lets users play a game of chess against the engine. It will accept a FEN string and a list of moves, and return the top N moves along with their evaluations.
    