package database

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func UpdatePost(id int, title string, content string) error {
	// Open database
	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare SQL query to UPDATE
	query, err := db.Prepare("UPDATE posts SET title=?, content=? WHERE id=?")
    if err != nil {
        return err
    }

	// Execute query to UPDATE
    _, err = query.Exec(title, content, id)
    if err != nil {
        return err
    }

	return nil
}
