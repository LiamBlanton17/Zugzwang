package backend

import (
	"fmt"
	"net/http"
	"zugzwang/internal/platform"

	"github.com/gin-gonic/gin"
)

// Struct to bind json POST data
type Setup struct {
	Name string `json:"name" binding:"required"`
	Elo  int    `json:"elo" binding:"required"`
}

/*
 *	route: "/setup"
 *  method: "POST"
 *	description: Called by frontend to request to setup a game
 *	returns: Either an error that the game could not be setup, or success with a game id
 */
func HandleSetup(c *gin.Context) {
	ctx := c.Request.Context()

	// Validate the request
	// Request must have a name/username of the user and their estimated elo
	// This is just for analytics, not security or anything
	var setup Setup
	if err := c.ShouldBindBodyWithJSON(&setup); err != nil {

		// Return bad request 400 error if request does not validate
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new game
	// TODO: There are legitmate reasons why a game my fail other than a server error (like too many current games running)
	gameId, err := platform.CreateGame(setup.Name, setup.Elo, ctx)
	if err != nil {
		fmt.Println("Failed to create a new game.")

		// Return generic error to user
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	// Return the game id to the user, this is what the frontend will use to establish the websocket connection later
	c.JSON(http.StatusCreated, gin.H{"gameId": gameId})
}
