package main

import (
	"fmt"
	"log"
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

		showAllPostsPage.Execute(w, database.SelectPosts())
	})

	// Update

	// Delete



	
	fmt.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
