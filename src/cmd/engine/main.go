package main

import (
	"encoding/json"
	"net/http"
	"zugzwang/engine"

	"flag"
	"fmt"
)

/*
This is the main binary for the command line interface with the chess engine
*/
func main() {

	var action string
	flag.StringVar(&action, "action", "perft", "the action the program takes")
	flag.Parse()

	switch action {
	case "perft":
		engine.Perft()
	case "strengthtest":
		engine.StrengthTest()
	case "benchmark":
		engine.RunBenchmark()
	case "api":
		RunAPI()
	default:
		fmt.Println("The action is not supported: ", action)
	}
}

func RunAPI() {
	// Setup the HTTP endpoint
	// This is a single endpoint that accepts a POST request and returns the best move(s)
	http.HandleFunc("/move", func(w http.ResponseWriter, r *http.Request) {

		// Setup the engine
		engine.InitEngine()

		// Require a POST request
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body for the PGN string and the number of moves to return
		var request struct {
			PGN   string `json:"pgn"`
			Moves int    `json:"moves"`
		}

		// Decode the JSON request body
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Get the best move(s) from the engine
		bestMoves := engine.GetBestMoves(request.PGN, request.Moves)

		// Encode the best move(s) as JSON and return it in the response
		response := struct {
			BestMoves []string `json:"best_moves"`
		}{
			BestMoves: bestMoves,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the HTTP server
	fmt.Println("Starting API server on :8080")
	http.ListenAndServe(":8080", nil)
}
