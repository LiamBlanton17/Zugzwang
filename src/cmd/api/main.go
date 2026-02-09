package main

import (
	"backend/internal/api"
	"backend/internal/platform"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting backend API.")

	// Setup the DB
	err := platform.InitDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	// -- Setup the http handlers

	// Setup gin router
	r := gin.Default()

	/*
	 *	route: "/setup"
	 *  method: "POST"
	 *	description: Called by the frontend to request to setup a game
	 *	returns: Either an error that the game could not be setup, or success with a game id
	 */
	r.POST("/setup", api.HandleSetup)

	/*
	 *	route: "/start/{game_id}"
	 *	method: "GET"
	 *	desrciption: Called by the frontend to request to start a game
	 *  returns: Either an error that the game could not be started, or upgrades to a WS connection
	 */
	r.GET("/start/:game_id", api.HandleGame)

	// Run the router
	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
  
}
