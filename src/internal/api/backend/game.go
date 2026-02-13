package backend

import (
	"cmp"
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"zugzwang/internal/engine"
	"zugzwang/internal/platform"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/*
 *	route: "/start/{game_id}"
 *	method: "GET"
 *	desrciption: Called by the frontend to request to started a game, upgrades
 *  returns: Either an error that the game could not be started, or upgrades to a WS connection
 */
func HandleGame(c *gin.Context) {
	ctx := c.Request.Context()

	// Get the game id out of the route
	gameId := c.Param("game_id")

	// Validate te game is pending in the database and start it
	err := platform.StartGame(gameId, ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Game ID does not exist")
		}

		fmt.Println("Failed to start a new game.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not start the game"})
		return
	}

	// Get the web socket upgrader and upgrade the http request
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {

		fmt.Println("Failed to upgrade protocol to websocket.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not start the game"})
		return
	}

	// Start the game
	gameLoop(ws)

	// Update the game in the database to be finished
}

// Game state struct
type gameState struct {
	board       *engine.Board
	playerColor engine.Color
	engineColor engine.Color
}

// Function to perform the game loop with the websocket
func gameLoop(ws *websocket.Conn) {
	defer ws.Close()

	// Setup a new game state
	state := gameState{}
	board, _ := engine.STARTING_POSITION_FEN.ToBoard(nil)
	state.board = board

	// Get the colors for the game
	playerColor, engineColor, err := determineGameColors(ws)
	if err != nil {
		fmt.Println(err)
		return
	}
	state.playerColor = playerColor
	state.engineColor = engineColor

	// Setting depth to just 7 ply as a default for not, likely refactor to be time constrained
	const depth = 7

	// Start game loop
	for {

		// Engines turn
		if state.engineColor == state.board.Turn {
			result := state.board.RootSearch(depth, false) // false heres means don't multithread

			// Get the moves and sort them
			moveResults := result.Moves
			slices.SortFunc(moveResults, func(a engine.MoveEval, b engine.MoveEval) int {
				return cmp.Compare(b.Eval, a.Eval)
			})

			// Eval needs to be context aware
			bestEval := moveResults[0].Eval
			if board.Turn == engine.BLACK {
				bestEval *= -1
			}
			bestMove := moveResults[0].Move

			// Todo: in future select the second or third best move based on some math equation

			// Make the move (this handles everything, including swapping the board colors)
			state.board.Move(bestMove)

			// Send that move over the web socket
			ws.WriteMessage(websocket.TextMessage, bestMove.ToBytes())

		}

		// Players turn
		if state.playerColor == state.board.Turn {

			// Get the legal moves in the position, and send them to the client
			moves := state.board.GenerateLegalMoves()

			// Send those to the client
			for _, move := range moves {
				ws.WriteMessage(websocket.TextMessage, move.ToBytes())
			}

			// Read their move in
			_, msg, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("Websocket error, trying to read the clients move")
				return
			}

			// Message is expected to be a length of 1 byte, corresponding to the index in the list of moves sent
			if len(msg) != 1 {
				fmt.Println("Websocket error, invalid move sent by client")
				return
			}

			moveIdx := int(msg[0])
			if moveIdx >= len(moves) {
				fmt.Println("Websocket error, invalid move sent by client")
				return
			}
			playerMove := moves[moveIdx]

			// Make the move
			state.board.Move(playerMove)
		}

	}

}

// Function to set the colors for each player
// First color is for the player, second is for the engine
func determineGameColors(ws *websocket.Conn) (engine.Color, engine.Color, error) {
	// Expecting the first byte to be equal to 0, 1, or 2 (engine.WHITE, engine.BLACK, engine.EITHER_COLOR)
	_, msg, err := ws.ReadMessage()
	if err != nil {
		fmt.Println("Failed to get the color of the player")
		return 0, 0, fmt.Errorf("Failed to get the color of the player because of an invalid selection")
	}

	// Require the message to be only 1 byte
	if len(msg) != 1 {
		fmt.Println("Invalid message read from the color of the player")
		return 0, 0, fmt.Errorf("Failed to get the color of the player because of an invalid selection")
	}

	// Make the color selection
	playerColor := engine.Color(msg[0])
	switch playerColor {
	case engine.WHITE:
		return engine.WHITE, engine.BLACK, nil
	case engine.BLACK:
		return engine.BLACK, engine.WHITE, nil
	case engine.EITHER_COLOR:
		// If player gives the choice to the engine, the engine will chose to be black
		// Refactor later to make a random choice
		return engine.WHITE, engine.BLACK, nil
	default:
		// If invalid message, return error
		return 0, 0, fmt.Errorf("Failed to get the color of the player because of an invalid selection")
	}

}
