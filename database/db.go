package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("Error opening database: %v", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
	}

	fmt.Println("Connected to PostgreSQL database!")
	return db
}

func GetDB() *sql.DB {
	return db
}
