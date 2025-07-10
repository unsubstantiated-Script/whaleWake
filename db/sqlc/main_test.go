package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain Entry point for all the tests that need to run.
func TestMain(m *testing.M) {
	var err error

	err = godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Utilizing .env vars, so keeping these local variables. b
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PWORD")
	dbDriver := "postgres"
	dbSource := fmt.Sprintf("postgresql://%s:%s@localhost:5432/whale_wake_users?sslmode=disable", dbUser, dbPassword)

	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Unable to connect to the db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
