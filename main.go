package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"whaleWake/api"
	db "whaleWake/db/sqlc"
	"whaleWake/util"
)

// main is the entry point of the application.
// It initializes environment variables, connects to the database,
// creates a new server instance, and starts the HTTP server.
func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Unable to load config:", err)
	}

	// Open a connection to the database.
	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("Unable to connect to the db:", err)
	}

	// Create a new store instance for database operations.
	store := db.NewStore(conn)

	// Create a new server instance with the store.
	server, err := api.NewServer(store)

	if err != nil {
		log.Fatal("Unable to start the server:", err)
	}

	// Start the HTTP server on the specified address.
	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("Unable to start the server:", err)
	}
}
