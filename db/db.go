package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func New() *sql.DB {
	return NewDB()
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./entries.sqlite3")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
