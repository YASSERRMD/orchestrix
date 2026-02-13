package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbPath string) error {
	if dbPath == "" {
		dbPath = "./orchestrix.db"
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func Migrate(migrationsPath string) error {
	migrations, err := os.ReadFile(migrationsPath)
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(migrations))
	return err
}
