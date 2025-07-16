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

func main() {
	var err error

	err = godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Utilizing .env vars, so keeping these local variables.
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PWORD")
	dbAddress := os.Getenv("DB_PORT_ADDY")
	dbDriver := "postgres"
	dbSource := fmt.Sprintf("postgresql://%s:%s@localhost:5432/whale_wake_users?sslmode=disable", dbUser, dbPassword)

	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Unable to connect to the db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(dbAddress)

	if err != nil {
		log.Fatal("Unable to start the server:", err)
	}
}
