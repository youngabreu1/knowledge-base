package database

import (
	"fmt"
	"database/sql"
	"github.com/joho/godotenv"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewConnection() (*sql.DB, error) {
	godotenv.Load()
	var dbProvider = os.Getenv("DB_PROVIDER")
	var dbUser = os.Getenv("DB_USER")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbHost = os.Getenv("DB_HOST")
	var dbName = os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s",dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open(dbProvider, connStr)
	if err != nil {
		return nil, fmt.Errorf("Error opening connection: %s", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error connecting to database: %w", err)
	}

	return db, nil
}
