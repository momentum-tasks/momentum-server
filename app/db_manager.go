package app

import (
	"database/sql"
	"time"

	// Import the driver for the mysql database
	_ "github.com/go-sql-driver/mysql"
)

var store *DBManager

// DBManager stores the sql.DB object for easier access throughout the package
type DBManager struct {
	db *sql.DB
}

// NewDBManager creates the DBManager that will be used to connect and post queries to the database
func NewDBManager() *DBManager {
	return new(DBManager)
}

// Begin opens the database connection, and attempts to connect to it
// If the connection is successful, then the database tables will be created if they do not already exist
func (db *DBManager) Begin(dbType string, connString string) {
	var err error
	db.db, err = sql.Open(dbType, connString)
	if err != nil {
		mlog.Fatal(err)
	} else {
		mlog.Info("Database connection open... Trying to ping...")
	}

	for i := 0; i < 3000; i++ {
		time.Sleep(100 * time.Millisecond)

		err = db.db.Ping()
		if err == nil {
			mlog.Info("Database connected successfully")
			break
		}
	}
	if err != nil {
		mlog.Fatal(err)
	}

	store = db

	createTables()
}

func createTables() error {
	_, err := store.db.Exec(userTableStatement)
	if err != nil {
		mlog.Fatal(err)
		return err
	}
	_, err = store.db.Exec(userSessionTableStatement)
	if err != nil {
		mlog.Fatal(err)
		return err
	}
	_, err = store.db.Exec(taskTableStatement)
	if err != nil {
		mlog.Fatal(err)
		return err
	}
	_, err = store.db.Exec(reportTableStatement)
	if err != nil {
		mlog.Fatal(err)
		return err
	}

	return nil
}
