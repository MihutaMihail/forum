package main

import (
	"fmt"
	"log"
	"strconv"
	"net/http"
	"html/template"
	"forum/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePage := template.Must(template.ParseFiles("../templates/index.html"))
		homePage.Execute(w, "o")
	})

	//
	// POSTS
	//

	// Create
	http.HandleFunc("/createPost", func(w http.ResponseWriter, r *http.Request) {
		createPostPage := template.Must(template.ParseFiles("../templates/posts/createPost.html"))

		createPostPage.Execute(w, "o")
	})

	http.HandleFunc("/submitPost", func(w http.ResponseWriter, r *http.Request) {
		database.InsertPost(r.FormValue("title"), r.FormValue("content"))

		http.Redirect(w, r, "/showAllPosts", http.StatusFound)
	})

	// Read
	http.HandleFunc("/showAllPosts", func(w http.ResponseWriter, r *http.Request) {
		showAllPostsPage := template.Must(template.ParseFiles("../templates/posts/showAllPosts.html"))

		showAllPostsPage.Execute(w, database.SelectAllPosts())
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the query parameters
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		// Change to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Retrieve the post from the database using the ID
		post := database.GetPostByID(id)
		if post == (database.Post{}) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		showPostPage := template.Must(template.ParseFiles("../templates/posts/showPost.html"))

		showPostPage.Execute(w, post)
	})
	
	fmt.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
