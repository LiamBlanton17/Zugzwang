package main

import (
	"fmt"
	"zugzwang/internal/api/backend"
	"zugzwang/internal/api/frontend"
	"zugzwang/internal/platform"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting application.")

	// Setup the DB
	err := platform.InitDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	// -- Setup the http handlers

	// Setup gin router
	r := gin.Default()

	// Setup static file handling
	r.Static("/static", "./static")

	/*
	 * route: "/"
	 * method: "GET"
	 * desrciption: the main frontend endpoint
	 * returns: HTML for the main layout
	 */
	r.GET("/", frontend.HandleIndex)

	/*
	 * route: "/learn-more"
	 * method: "GET"
	 * desrciption: the learn more page frontend endpoint
	 * returns: HTML for the learn more
	 */
	r.GET("/learn-more", frontend.HandleLearnMore)

	/*
	 * route: "/game"
	 * method: "GET"
	 * description: the game page frontend endpoint
	 * returns: HTML for the game UI
	 */
	r.GET("/game", frontend.HandleGame)

	/*
	 *	route: "/setup"
	 *  method: "POST"
	 *	description: Called by the frontend to request to setup a game
	 *	returns: Either an error that the game could not be setup, or success with a game id
	 */
	r.POST("/setup", backend.HandleSetup)

	/*
	 *	route: "/start/{game_id}"
	 *	method: "GET"
	 *	desrciption: Called by the frontend to request to start a game
	 *  returns: Either an error that the game could not be started, or upgrades to a WS connection
	 */
	r.GET("/start/:game_id", backend.HandleGame)

	// Run the router
	if err := r.Run(":8067"); err != nil {
		fmt.Println(err)
	}

}
