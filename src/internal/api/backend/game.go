package backend

import (
	"database/sql"
	"fmt"
	"net/http"
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

		fmt.Println("Failed to ypgrade protocol to websocket.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not start the game"})
		return
	}

	// Start the game
	gameLoop(ws, c)
}

// Function to perform the game loop with the websocket
func gameLoop(ws *websocket.Conn, c *gin.Context) {
	defer ws.Close()
}
