package database

import (
	"database/sql"
	"flowforge/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

var DB *sql.DB

func Initialize(cfg *config.Config) error {
	dbPath := cfg.DBPath

	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=ON")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
