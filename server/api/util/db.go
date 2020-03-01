package util

import (
	"database/sql"
	"errors"
	"fmt"
)

var db *sql.DB

// GetDB method returns a DB instance
func GetDB() (*sql.DB, error) {
	// connectionString := "user=postgres password=sai dbname=postgres sslmode=disable search_path=saiteja"
	connectionString := "postgres://cmhqeyjitcetwh:69b345b70185e115bc9f6722a6389d6c282cdab4fa53aeb01192132c24ebc3d8@ec2-184-72-236-3.compute-1.amazonaws.com:5432/degm90gbclc7ps"
	if connectionString == "" {
		return nil, errors.New("'POSTGRES_CONNECTION_STRING' environment variable not set")
	}
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}
	return conn, nil
}
