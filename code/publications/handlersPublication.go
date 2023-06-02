package publications

import (
	"bytes"
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

type indexPageData struct { // maybe temporary ; interface template with main body
	Main   template.HTML
	Header template.HTML
}

type mainFeedData struct {
	Publications []template.HTML
}

type headerData struct {
	HeaderData      template.HTML
	FilterBarOn 	bool
	IsUserConnected bool
}

//
// READ
//

/*
Write main feed
*/
func HandleAllPosts(w http.ResponseWriter, r *http.Request) {
	indexData := indexPageData{}

	// get tag value (filter)
	r.ParseForm()
	tag := r.FormValue("tag")

	// make mainFeed
	mainFeed := mainFeedData{}
	mainFeed.Publications = SortAllPublication(w, r, tag)

	tplMain := new(bytes.Buffer)
	tplRawMain := template.Must(template.ParseFiles("templates/mainFeed.html"))
	err := tplRawMain.Execute(tplMain, mainFeed)
	checkErr(err)
	tplStringMain := tplMain.String()
	indexData.Main = template.HTML(tplStringMain)

	indexData.Header = MakeHeaderTemplate(w, r, true)

	// execute with interface
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

	// Retrieve tags
	post.TagsString = GetTagsString(post.Pid)

	// Change to UPDATE POST
	formPost.Execute(w, post)
}

func HandleSubmitForm(w http.ResponseWriter, r *http.Request) {
	var newPost PublicationData

	r.ParseForm()

	//
	// Image
	//

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
	if file != nil && r.FormValue("addImageBoolean") == "true" {
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
	} else if (r.FormValue("imageName") != "") && r.FormValue("addImageBoolean") == "true" {
		newPost.ImageLink = r.FormValue("imageName")
	}

	//
	// Title, Content, Tags
	//

	newPost.Title = r.FormValue("title")
	newPost.Content = r.FormValue("content")
	if file != nil && r.FormValue("addImageBoolean") == "true" {
		newPost.ImageLink = filename
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

	if r.FormValue("pid") == "-1" {
		err := InsertPost(newPost, selectedTags, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// UPDATE POST
		postId, err := strconv.Atoi(r.FormValue("pid"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		post := GetPostByID(postId)
		post.Title = newPost.Title
		post.Content = newPost.Content
		post.ImageLink = newPost.ImageLink
		post.Edited = 1

		err = UpdatePost(post, selectedTags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

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
	post := GetPostByID(id)

	sessionUser := authentification.GetSessionUid(w, r)
	isAdmin := isUserAdmin(w, r)

	uidInt, err := strconv.Atoi(post.Uid)
	checkErr(err)
	// check if the user is the owner, or an admin
	if uidInt == sessionUser || isAdmin {
		DeletePost(post)
	}
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

// return true if sessionUser is admin, else false
func isUserAdmin(w http.ResponseWriter, r *http.Request) bool {
	// open Db
	sessionUser := authentification.GetSessionUid(w, r)
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// if not connected
	if sessionUser == 0 {
		return false
	}

	preparedRequest, err := db.Prepare("SELECT admin FROM Users WHERE uid = ?;")
	checkErr(err)
	isAdmin := 0
	row := preparedRequest.QueryRow(sessionUser)
	row.Scan(&isAdmin)

	return isAdmin == 1
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// get the cid from query url
	cid := r.URL.Query().Get("id")

	// get the pid of the comm
	preparedRequest, err := db.Prepare("SELECT pid FROM Comments WHERE cid = ?;")
	checkErr(err)
	pid := 0
	row := preparedRequest.QueryRow(cid)
	row.Scan(&pid)

	// delete the comm
	_, err = db.Exec("DELETE FROM Comments WHERE cid = ?;", cid)
	checkErr(err)

	refreshPublicationPage(w, r, pid)
}
