package publications

import (
	"database/sql"
	"forum/code/authentification"
	"log"
	"math"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	posts         []PublicationData
	existingPosts = make(map[int]bool)
)

//
// SELECT
//

func GetAllPosts(postValue string, uid int) []PublicationData {
	// Open database
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	var query string
	var args []interface{}

	if postValue == "myPosts" {
		query = "SELECT * FROM publications WHERE uid=?"
		args = []interface{}{uid}
	} else if postValue == "likedPosts" {
		query = "SELECT p.* FROM Publications p INNER JOIN Likes l ON p.pid = l.pid WHERE l.uid=? AND l.isLike = 1;"
		args = []interface{}{uid}
	} else {
		query = "SELECT * FROM publications"
		args = nil
	}

	rows, err := db.Query(query, args...)
	checkErr(err)
	defer rows.Close()


	posts = nil
	existingPosts = make(map[int]bool)
	// Store the select posts
	for rows.Next() {
		var post PublicationData
		err := rows.Scan(
			&post.Pid, &post.Title, &post.Content, &post.ImageLink,
			&post.UpvoteNumber, &post.CreatedDate, &post.Edited, &post.Uid)
		checkErr(err)

		// RATINGS
		timeNow := time.Now().Format("02-01-2006")
		timeStart, err := time.Parse("02/01/2006", post.CreatedDate)
		checkErr(err)
		timeEnd, err := time.Parse("02-01-2006", timeNow)
		checkErr(err)
		days := math.Ceil(timeEnd.Sub(timeStart).Hours() / 24)
		post.Rating = post.UpvoteNumber + post.CommentNumber - int(math.Round(math.Pow(days, 2)))

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

func GetTagsString(pid int) []string {
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	preparedRequest, err := db.Prepare("SELECT name FROM Tags WHERE pid = ?;")
	checkErr(err)
	rows, err := preparedRequest.Query(pid)
	checkErr(err)
	defer rows.Close()

	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		checkErr(err)
		tagArray = append(tagArray, tag)
	}

	return tagArray
}

//
// INSERT
//

func InsertPost(post PublicationData, selectedTags []string, w http.ResponseWriter, r *http.Request) error {
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	query, err := db.Prepare("INSERT INTO publications(title, content, image, like, createdDate, uid) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer query.Close()

	// Execute query to INSERT
	_, err = query.Exec(post.Title, post.Content, post.ImageLink, 0, time.Now().Format("02/01/2006"), authentification.GetSessionUid(w, r))
	if err != nil {
		return err
	}

	// Get last POST which is current POST
	posts := GetAllPosts("", 0)
	lastPost := posts[len(posts)-1]

	// Insert tags
	for _, tag := range selectedTags {
		query, err := db.Prepare("INSERT INTO tags(name, pid) VALUES(?, ?)")
		if err != nil {
			return err
		}
		defer query.Close()

		_, err = query.Exec(tag, lastPost.Pid)
		if err != nil {
			return err
		}
	}

	return nil
}

//
// DELETE
//

func DeletePost(post PublicationData) error {
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	query, err := db.Prepare("DELETE FROM publications WHERE pid=?")
	if err != nil {
		return err
	}
	defer query.Close()
	_, err = query.Exec(post.Pid)
	if err != nil {
		return err
	}

	query, err = db.Prepare("DELETE FROM tags WHERE pid=?")
	if err != nil {
		return err
	}
	defer query.Close()
	_, err = query.Exec(post.Pid)
	if err != nil {
		return err
	}

	query, err = db.Prepare("DELETE FROM Comments WHERE pid=?")
	if err != nil {
		return err
	}
	defer query.Close()
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

func UpdatePost(post PublicationData, selectedTags []string) error {
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// POST
	query, err := db.Prepare("UPDATE publications SET title=?, content=?, image=?, edited=? WHERE pid=?")
	if err != nil {
		return err
	}

	_, err = query.Exec(post.Title, post.Content, post.ImageLink, post.Edited, post.Pid)
	if err != nil {
		return err
	}

	// TAGS
	query, err = db.Prepare("DELETE from tags WHERE pid=?")
	if err != nil {
		return err
	}

	_, err = query.Exec(post.Pid)
	if err != nil {
		return err
	}

	for _, tag := range selectedTags {
		query, err := db.Prepare("INSERT INTO tags(name, pid) VALUES(?, ?)")
		if err != nil {
			return err
		}
		defer query.Close()

		_, err = query.Exec(tag, post.Pid)
		if err != nil {
			return err
		}
	}

	return nil
}
