package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:Aquamarine@localhost:5432/aquatrade?sslmode=disable"
)

func TestMain(m *testing.M) {
	testDB, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
