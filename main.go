package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"whaleWake/api"
	db "whaleWake/db/sqlc"
)

// main is the entry point of the application.
// It initializes environment variables, connects to the database,
// creates a new server instance, and starts the HTTP server.
func main() {
	var err error

	// Load environment variables from the .env file.
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve database credentials and address from environment variables.
	dbUser := os.Getenv("DB_USER")         // Database username.
	dbPassword := os.Getenv("DB_PWORD")    // Database password.
	dbAddress := os.Getenv("DB_PORT_ADDY") // Server address for the database.
	dbDriver := "postgres"                 // Database driver (PostgreSQL).
	dbSource := fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432/whale_wake_users?sslmode=disable",
		dbUser, dbPassword, // Connection string for the database.
	)

	// Open a connection to the database.
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Unable to connect to the db:", err)
	}

	// Create a new store instance for database operations.
	store := db.NewStore(conn)

	// Create a new server instance with the store.
	server := api.NewServer(store)

	// Start the HTTP server on the specified address.
	err = server.Start(dbAddress)
	if err != nil {
		log.Fatal("Unable to start the server:", err)
	}
}
