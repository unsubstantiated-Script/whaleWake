package db

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
	"whaleWake/util"
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

	// Load environment variables from .env file
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Unable to connect to the db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
