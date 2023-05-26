package publications

import (
	"database/sql"
	"encoding/json"
	"forum/code/authentification"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type indexPageData struct {
	Publications []template.HTML
}

//
// READ
//

func HandleAllPosts(w http.ResponseWriter, r *http.Request) {
	indexData := indexPageData{}
	for _, post := range GetAllPosts() {
		publication := MakePublicationHomePageTemplate(post.Pid, w, r)
		indexData.Publications = append(indexData.Publications, publication)
	}

	allPosts := template.Must(template.ParseFiles("./templates/publicationListTemplate.html"))
	allPosts.Execute(w, indexData)
}

//
// CREATE
//

func CheckHandleFormPost(w http.ResponseWriter, r *http.Request) {
	if authentification.CheckSessionUid(w, r) == nil {
		http.Redirect(w, r, "/publicationForm", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", http.StatusNotFound)

	}
}

func HandleFormPost(w http.ResponseWriter, r *http.Request) {
	formPost := template.Must(template.ParseFiles("./templates/publicationFormTemplate.html"))
	// Get the ID from the query parameters
	idStr := r.URL.Query().Get("id")

	// Change to CREATE POST
	if idStr == "" {
		var postEmpty PublicationData
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
	post := GetPostByID(id)
	if post.Title == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Database - Get Tags
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	preparedRequest, err := db.Prepare("SELECT name FROM Tags WHERE pid = ?;")
	checkErr(err)
	rows, err := preparedRequest.Query(post.Pid)
	checkErr(err)
	defer rows.Close()

	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		checkErr(err)
		tagArray = append(tagArray, tag)
	}
	post.TagsString = tagArray

	// Change to UPDATE POST
	formPost.Execute(w, post)
}

func HandleSubmitForm(w http.ResponseWriter, r *http.Request) {
	var post PublicationData

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
		out, err := os.Create("./assets/uploads/" + filename)
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

	r.ParseForm()

	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	if file != nil {
		post.ImageLink = filename
	}
	selectedTagsJSON := r.FormValue("selected-tags")
	var selectedTags []string

	if selectedTagsJSON != "" {
		err = json.Unmarshal([]byte(selectedTagsJSON), &selectedTags)
		checkErr(err)

		for i, tag := range selectedTags {
			selectedTags[i] = strings.TrimSuffix(tag, "x")
		}
	}

	InsertPost(post, selectedTags)

	http.Redirect(w, r, "/", http.StatusFound)
}

//
// DELETE
//

func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := getQueryID(w, r)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	DeletePost(GetPostByID(id))
	http.Redirect(w, r, "/", http.StatusFound)
}

//
// General Functions
//

func getQueryID(w http.ResponseWriter, r *http.Request) (int, error) {
	// Get the ID from the query parameters
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}
	return id, err
}
