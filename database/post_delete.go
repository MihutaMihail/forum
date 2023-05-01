package database

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func DeletePost(id int) error {
	// Open database
	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare SQL query to INSERT into POSTS
	query, err := db.Prepare("DELETE FROM posts WHERE id=?")
    if err != nil {
        return err
    }
    defer query.Close()

	// Execute query to INSERT
	_, err = query.Exec(id)
    if err != nil {
		return err
    }

	deleteFromArray(id)

	return nil
}

func deleteFromArray(id int) {
	postDelete := GetPostByID(id)

	for i, post := range posts {
		if post.ID == postDelete.ID {
			posts = append(posts[:i], posts[i+1:]...)
            break
		}
	}
}