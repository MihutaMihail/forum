package publications

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	posts         []PublicationData
	existingPosts = make(map[int]bool)
)

//
// SELECT
//

func GetAllPosts() []PublicationData {
	// Open database
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// Prepare SQL query to SELECT ALL POSTS
	rows, err := db.Query("SELECT * FROM publications")
	checkErr(err)
	defer rows.Close()

	// Store the select posts
	for rows.Next() {
		var post PublicationData
		err := rows.Scan(
			&post.Pid, &post.Title, &post.Content, &post.ImageLink,
			&post.UpvoteNumber, &post.CreatedDate, &post.Uid)
		checkErr(err)

		// Check if post exists
		if _, ok := existingPosts[post.Pid]; !ok {
			posts = append(posts, post)
		} else {
			// Update existing posts
			for i := range posts {
				if posts[i].Pid == post.Pid {
					posts[i] = post
					break
				}
			}
		}
		existingPosts[post.Pid] = true
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return posts
}

func GetPostByID(id int) PublicationData {
	for _, post := range posts {
		if post.Pid == id {
			return post
		}
	}
	return PublicationData{}
}

//
// INSERT
//

func InsertPost(post PublicationData) error {
	// Open database
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// Prepare SQL query to INSERT into POSTS
	query, err := db.Prepare("INSERT INTO publications(title, content, image, like, createdDate, uid) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer query.Close()

	// Execute query to INSERT
	_, err = query.Exec(post.Title, post.Content, post.ImageLink, post.UpvoteNumber, post.CreatedDate, post.Uid)
	if err != nil {
		return err
	}

	return nil
}

//
// DELETE
//

func DeletePost(post PublicationData) error {
	// Open database
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// Prepare SQL query to INSERT into POSTS
	query, err := db.Prepare("DELETE FROM publications WHERE pid=?")
	if err != nil {
		return err
	}
	defer query.Close()

	// Execute query to INSERT
	_, err = query.Exec(post.Pid)
	if err != nil {
		return err
	}

	deleteFromArray(post.Pid)

	return nil
}

func deleteFromArray(id int) {
	postDelete := GetPostByID(id)
	for i, post := range posts {
		if post.Pid == postDelete.Pid {
			posts = append(posts[:i], posts[i+1:]...)
			break
		}
	}
}

//
// UPDATE
//

func UpdatePost(post PublicationData) error {
	// Open database
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// Prepare SQL query to UPDATE
	query, err := db.Prepare("UPDATE publications SET title=?, content=?, image=?, like=? WHERE pid=?")
	if err != nil {
		return err
	}

	// Execute query to UPDATE
	_, err = query.Exec(post.Title, post.Content, post.ImageLink, post.UpvoteNumber)
	if err != nil {
		return err
	}

	return nil
}
