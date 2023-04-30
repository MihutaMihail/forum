package database

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InsertPost(title string, content string) error {
	// Open database
	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare SQL query to INSERT into POSTS
	query, err := db.Prepare("INSERT INTO posts(title, content) VALUES(?, ?)")
    if err != nil {
        return err
    }
    defer query.Close()

	// Execute query to INSERT
	_, err = query.Exec(title, content)
    if err != nil {
		return err
    }

	return nil
}