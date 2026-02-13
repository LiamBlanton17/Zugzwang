package backend

import (
	"cmp"
	"context"
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
	gameLoop(ctx, ws)

	// Update the game in the database to be finished
}

// Game message types
type gameMessageType byte

const (
	PLAYER_COLOR_SELECT gameMessageType = iota
	ENGINE_COLOR_CONFIRM
	PLAYER_MAKE_MOVE
	ENGINE_MAKE_MOVE
	LIST_LEGAL_MOVES
	GAME_DRAW
	GAME_ENGINE_WIN
	GAME_PLAYER_WIN
	GAME_QUIT
	DRAW_OFFER
	DRAW_REJECT
	DRAW_ACCEPT
)

// Game message struct
type gameMessage struct {
	gmType        gameMessageType
	payloadLength int
	payload       []byte
}

// Function to write the game message
// Format is [Type][Payload Length][Payload]
func (msg *gameMessage) write(ws *websocket.Conn) error {
	payloadLen := len(msg.payload)
	message := make([]byte, 0, 3+payloadLen)
	message = append(message, byte(msg.gmType))
	message = append(message, byte(payloadLen>>8), byte(payloadLen&0xFF))
	message = append(message, msg.payload...)
	return ws.WriteMessage(websocket.BinaryMessage, message)
}

// Function to read a message from the client into a game message struct
func (msg *gameMessage) read(ws *websocket.Conn) error {
	mt, message, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	// Make sure it is binary message
	if mt != websocket.BinaryMessage {
		return fmt.Errorf("Websocket did not send a binary message")
	}

	// Make sure the message is at least 3 bytes
	msg.payloadLength = len(message)
	if msg.payloadLength < 3 {
		return fmt.Errorf("Websocket message must be at least 3 bytes (type byte, 2 bytes for length)")
	}

	// Handle the message
	msg.gmType = gameMessageType(message[0])
	if msg.payloadLength > 3 {
		msg.payload = message[3:]
	} else {
		msg.payload = nil
	}

	return nil
}

// Game state struct
type gameState struct {
	board       *engine.Board
	playerColor engine.Color
	engineColor engine.Color
}

// Function to perform the game loop with the websocket
func gameLoop(ctx context.Context, ws *websocket.Conn) {
	defer ws.Close()

	// Setup a new game state
	state := gameState{}
	board, _ := engine.STARTING_POSITION_FEN.ToBoard(nil)
	state.board = board

	// Create a message struct to read/write all messages to and from
	msg := &gameMessage{}

	// Get the colors for the game
	playerColor, engineColor, err := determineGameColors(ws, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	state.playerColor = playerColor
	state.engineColor = engineColor

	// Send those back to the client for confirmation
	msg.gmType = ENGINE_COLOR_CONFIRM
	msg.payload = []byte{byte(playerColor), byte(engineColor)}
	if err := msg.write(ws); err != nil {
		fmt.Println(err)
		return
	}

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
			if state.board.Turn == engine.BLACK {
				bestEval *= -1
			}
			bestMove := moveResults[0].Move

			// Todo: in future select the second or third best move based on some math equation
			// Also, in the future the engine should be able to resign or offer draws

			// Make the move (this handles everything, including swapping the board colors)
			state.board.Move(bestMove)

			// Send that move over the web socket
			msg.gmType = ENGINE_MAKE_MOVE
			msg.payload = bestMove.ToBytes()
			if err := msg.write(ws); err != nil {
				fmt.Println(err)
				return
			}

		} else { // Players turn

			// Get the legal moves in the position, and send them to the client
			moves := state.board.GenerateLegalMoves()

			// Send those to the client
			msg.gmType = LIST_LEGAL_MOVES
			for _, move := range moves {
				msg.payload = append(msg.payload, move.ToBytes()...)
			}
			if err := msg.write(ws); err != nil {
				fmt.Println(err)
				return
			}

			// Read their move in
			if err := msg.read(ws); err != nil {
				fmt.Println("Websocket error, trying to read the clients move")
				return
			}

			// To do: refactor to allow the players to offer a draw

			// Message payload is expected to be 1 byte, corresponding to the index of the move to be played from the list of legal mvoes
			if msg.payloadLength != 1 {
				fmt.Println("Websocket error, trying to read the clients move, was not only 1 byte in size")
				return
			}

			// Try to get move from idx in moves
			moveIdx := int(msg.payload[0])
			if moveIdx >= len(moves) {
				fmt.Println("Websocket error, invalid move sent by client, out of range")
				return
			}
			playerMove := moves[moveIdx]

			// Make the move
			state.board.Move(playerMove)
		}

		// Handle game over conditions
	}

}

// Function to set the colors for each player
// First color is for the player, second is for the engine
func determineGameColors(ws *websocket.Conn, msg *gameMessage) (engine.Color, engine.Color, error) {

	// Attempt to read the player choice
	if err := msg.read(ws); err != nil {
		return 0, 0, err
	}

	// Require the message payload to be only 1 byte
	if msg.payloadLength != 1 {
		return 0, 0, fmt.Errorf("Failed to get the color of the player because of an invalid selection")
	}

	// Make the color selection
	playerColor := engine.Color(msg.payload[0])
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
