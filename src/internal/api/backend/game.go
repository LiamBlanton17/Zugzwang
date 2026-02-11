package backend

import "github.com/gin-gonic/gin"

/*
 *	route: "/start/{game_id}"
 *	method: "GET"
 *	desrciption: Called by the frontend to request to started a game, upgrades
 *  returns: Either an error that the game could not be started, or upgrades to a WS connection
 */
func HandleGame(c *gin.Context) {
	//gameId := c.Param("game_id")

}
