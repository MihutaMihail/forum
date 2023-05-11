package testcrud

import (
	"database/sql"
	"forum/code/publications"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type indexPageData struct {
	Publications []template.HTML
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	homePage := template.Must(template.ParseFiles("./code/testcrud/index.html"))
	homePage.Execute(w, "o")
}

func HandleAllPosts(w http.ResponseWriter, r *http.Request) {
	indexData := indexPageData{}
	for _, post := range publications.GetAllPosts() {
		publication := publications.MakePublicationHomePageTemplate(post.Pid)
		indexData.Publications = append(indexData.Publications, publication)
	}

	allPosts := template.Must(template.ParseFiles("./code/testcrud/allposts.html"))
	allPosts.Execute(w, indexData)
}

func HandleFormPost(w http.ResponseWriter, r *http.Request) {
	formPost := template.Must(template.ParseFiles("./code/testcrud/formPost.html"))
	// Get the ID from the query parameters
	idStr := r.URL.Query().Get("id")

	// Change to CREATE POST
	if idStr == "" {
		var postEmpty publications.PublicationData
		postEmpty.Pid = -1

		formPost.Execute(w, postEmpty)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	// Retrieve the post from the database using the ID
	post := publications.GetPostByID(id)
	if post.Title == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Change to UPDATE POST
	formPost.Execute(w, post)
}

func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := getQueryID(w, r)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	publications.DeletePost(publications.GetPostByID(id))
	http.Redirect(w, r, "/testAllPosts", http.StatusFound)
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	id, err := getQueryID(w, r)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	// Retrieve the post from the database using the ID
	post := publications.GetPostByID(id)
	if post.Title == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	showPost := template.Must(template.ParseFiles("./code/testcrud/post.html"))
	showPost.Execute(w, post)
}

func HandleSubmitForm(w http.ResponseWriter, r *http.Request) {
	var post publications.PublicationData

	file, header, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			file = nil
		} else {
			log.Println("Error retrieving file:", err)
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var filename string
	if file != nil {
		filename = header.Filename

		// Save the file
		out, err := os.Create("./code/testcrud/uploads/" + filename)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			log.Println("Error copying file:", err)
			http.Error(w, "Failed to copy file", http.StatusInternalServerError)
			return
		}
	}

	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	if file != nil {
		post.ImageLink = filename
	}

	db, err := sql.Open("sqlite3", "./database.db")
	defer db.Close()

	// Make tags
	preparedRequest, err := db.Prepare("SELECT name FROM Tags WHERE pid = ?")
	rows, err := preparedRequest.Query(post.Pid)
	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		tagArray = append(tagArray, tag)
	}
	post.Tags = publications.MakeTags(tagArray)

	publications.InsertPost(post)

	http.Redirect(w, r, "/testAllPosts", http.StatusFound)
}

func getQueryID(w http.ResponseWriter, r *http.Request) (int, error) {
	// Get the ID from the query parameters
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}
	return id, err
}
