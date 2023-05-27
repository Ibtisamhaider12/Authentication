//go:build cgo

package core

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
)

func init() {
	// Registers the sqlite3 driver with a ConnectHook so that we can
	// initialize the default PRAGMAs.
	//
	// Note 1: we don't define the PRAGMA as part of the dsn string
	// because not all pragmas are available.
	//
	// Note 2: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	sql.Register("sqlite3_with_connect_hook",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				conn.Exec(`
					PRAGMA busy_timeout       = 10000;
					PRAGMA journal_mode       = WAL;
					PRAGMA journal_size_limit = 200000000;
					PRAGMA synchronous        = NORMAL;
					PRAGMA foreign_keys       = ON;
				`, nil)

				return nil
			},
		},
	)
}

func connectDB(dbPath string) (*dbx.DB, error) {
	db, err := dbx.Open("sqlite3_with_connect_hook", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
