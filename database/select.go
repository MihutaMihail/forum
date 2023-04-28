package database

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	posts []Post
)

// Declaration of structures
type Post struct {
	ID int
	Title string
	Content string
}

func SelectPosts() []Post {
	posts = []Post{}

	// Open database
	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// CREATE TABLE POSTS
	// _, err = db.Exec("CREATE TABLE posts (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, content TEXT)")
	// if err != nil {
	// 	panic(err)
	// }

	// Prepare SQL query to INSERT into POSTS
	// query, err := db.Prepare("INSERT INTO posts(title, content) VALUES(?, ?)")
    // if err != nil {
    //     fmt.Println(err)
    // }
    // defer query.Close()

	// Execute query to INSERT
	// _, err = query.Exec("test2", "Content of test2")
    // if err != nil {
    //     fmt.Println(err)
    // }

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
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return posts
}