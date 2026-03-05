package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"github.com/google/uuid"
)

// This file holds the database functions for the project, which are very simplistic

// Global and Function used to initialize a global database pool
// Could refactor into DI if the project grows, but the project is simple enough as it is
// db will be not be exported from platform package, so it limits its scope
// Forces SQL to be written in this package, which helps package boundaries
var db *sql.DB

const sqlitePath = "./sqlite/sqlite.db"
const initDBPath = "./sqlite/init.sql"

func InitDB() error {
	// Attempt the connection
	tempdb, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return err
	}

	// Verify the connection
	if err := tempdb.Ping(); err != nil {
		return err
	}

	// Run init database
	sqlCode, err := os.ReadFile(initDBPath)
	if err != nil {
		return err
	}

	_, err = tempdb.Exec(string(sqlCode))
	if err != nil {
		return err
	}

	// Set the global and return no error
	db = tempdb
	return nil
}

// Function used to create a new game
// Validates we only have so many active games
const MAX_ACTIVE_GAMES = 100

func CreateGame(name string, elo int, ctx context.Context) (string, error) {

	var activeCount int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) AS ActiveGames FROM games WHERE status IN ('Active', 'Pending')").Scan(&activeCount)
	if err != nil {
		return "", err
	}

	if activeCount >= MAX_ACTIVE_GAMES {
		return "", fmt.Errorf("Too many games")
	}

	id := uuid.NewString()
	query := "INSERT INTO games (id, name, elo, status) VALUES (?, ?, ?, 'Pending')"
	_, err = db.ExecContext(ctx, query, id, name, elo)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return id, nil
}

// Function used to start a new game
func StartGame(gameId string, ctx context.Context) error {

	var activeCount int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) AS ActiveGames FROM games WHERE status IN ('Active', 'Pending')").Scan(&activeCount)
	if err != nil {
		return err
	}

	if activeCount >= MAX_ACTIVE_GAMES {
		return fmt.Errorf("Too many games")
	}

	err = db.QueryRowContext(ctx, "SELECT 1 FROM games WHERE status = 'Pending' AND id = ?", gameId).Scan()
	if err != nil {
		// Caller should check for sql.ErrNoRows too
		if err == sql.ErrNoRows {
			fmt.Println("Game ID does not exist")
		}
		return err
	}

	_, err = db.ExecContext(ctx, "UPDATE games SET status = 'Active' WHERE id = ?", gameId)
	if err != nil {
		return err
	}

	return nil
}
