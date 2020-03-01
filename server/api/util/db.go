package util

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

var db *sql.DB

// GetDB method returns a DB instance
func GetDB() (*sql.DB, error) {
	// connectionString := "user=postgres password=sai dbname=postgres sslmode=disable search_path=saiteja"
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return nil, errors.New("'POSTGRES_CONNECTION_STRING' environment variable not set")
	}
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}
	return conn, nil
}
