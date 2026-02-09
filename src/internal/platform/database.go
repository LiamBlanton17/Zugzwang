package platform

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// This file holds the database functions for the project, which are very simplistic

// Global and Function used to initialize a global database pool
// Could refactor into DI if the project grows, but the project is simple enough as it is
// db will be not be exported from platform package, so it limits its scope
// Forces SQL to be written in this package, which helps package boundaries
var db *sql.DB

func InitDB() error {
	// Attempt the connection
	tempdb, err := sql.Open("sqlite3", "./sqlite/sqlite.db")
	if err != nil {
		return err
	}

	// Verify the connection
	if err := tempdb.Ping(); err != nil {
		return err
	}

	// Set the global and return no error
	db = tempdb
	return nil
}

// Function used to create a new game
// Validates we only have so many active games
const MAX_ACTIVE_GAMES = 10

func CreateGame(name string, elo int, ctx context.Context) (string, error) {

	var activeCount int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) AS ActiveGames FROM games WHERE Status IN ('Active', 'Pending')").Scan(&activeCount)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	query := "INSERT INTO games (ID, Name, Elo, Status) VALUES (?, ?, ?, 'Pending')"
	_, err = db.ExecContext(ctx, query, id, name, elo)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return id, nil
}
