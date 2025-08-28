package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:8800/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Unable to establish connection with an error thrown : ", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
