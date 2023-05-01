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

	// CREATE & UPDATE
	http.HandleFunc("/formPost", func(w http.ResponseWriter, r *http.Request) {
		formPost := template.Must(template.ParseFiles("../templates/posts/formPost.html"))

		// Get the ID from the query parameters
		idStr := r.URL.Query().Get("id")

		// Change to CREATE POST
		if idStr == "" {
			var postEmpty database.Post
			postEmpty.ID = -1

			formPost.Execute(w, postEmpty)
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

		// Change to UPDATE POST
		formPost.Execute(w, post)
	})

	http.HandleFunc("/submitPost", func(w http.ResponseWriter, r *http.Request) {
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

		if (id == -1) {
			err := database.InsertPost(r.FormValue("title"), r.FormValue("content"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := database.UpdatePost(id, r.FormValue("title"), r.FormValue("content"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/showAllPosts", http.StatusFound)
	})

	// READ
	http.HandleFunc("/showAllPosts", func(w http.ResponseWriter, r *http.Request) {
		showAllPosts := template.Must(template.ParseFiles("../templates/posts/showAllPosts.html"))

		showAllPosts.Execute(w, database.SelectAllPosts())
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

		showPost := template.Must(template.ParseFiles("../templates/posts/showPost.html"))

		showPost.Execute(w, post)
	})
	
	// DELETE
	http.HandleFunc("/deletePost", func(w http.ResponseWriter, r *http.Request) {
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

		database.DeletePost(id)

		http.Redirect(w, r, "/showAllPosts", http.StatusFound)
	})

	fmt.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
