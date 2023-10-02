package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"codemead.com/go_fintech/fintech_backend/utils"
	_ "github.com/lib/pq"
)

var testQuery *db.Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}

	conn, err := sql.Open(config.DBdriver, config.DB_Source)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	testQuery = db.NewStore(conn)

	os.Exit(m.Run())
}
