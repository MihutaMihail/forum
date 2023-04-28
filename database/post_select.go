package database

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	posts []Post
	existingPosts = make(map[int]bool)
)

// Declaration of structures
type Post struct {
	ID int
	Title string
	Content string
}

func SelectAllPosts() []Post {
	// Open database
	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// CREATE TABLE POSTS IF NOT EXISTS
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, content TEXT)")
	if err != nil {
		panic(err)
	}

	// Prepare SQL query to SELECT ALL POSTS
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Store the select posts
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := existingPosts[post.ID]; !ok {
			posts = append(posts, post)
			existingPosts[post.ID] = true
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return posts
}

func GetPostByID(id int) Post {
	for _, post := range posts {
		if post.ID == id {
			return post
		}
	}
	return Post{}
}