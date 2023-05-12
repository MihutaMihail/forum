package publications

import (
	"database/sql"
	"forum/code/authentification"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

/*
Handle for the first click to add a comment ; call the same page with a commentBox
*/
func MakeCommentBox(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	checkErr(err)
	publicationData := makePublicationWithId(pid, w, r, "addCommentBox")

	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))

	err = tpl.Execute(w, publicationData)
	checkErr(err)
}

/*
Handle to add a comment in the final click
*/
func AddAComment(w http.ResponseWriter, r *http.Request) {
	var err error
	commentData := CommentData{}

	cookie := authentification.GetSessionUid(w, r)

	if cookie != 0 {
		commentData.Uid = cookie
		
		commentData.Pid, err = strconv.Atoi(r.URL.Query().Get("pid"))
		commentData.Like = 0

		timeNow := time.Now()
		commentData.CreatedDate = timeNow.Format("02-01-2006")

		commentData.Content = r.FormValue("content")

		// make a new Cid
		db, err := sql.Open("sqlite3", "./database.db")
		checkErr(err)
		defer db.Close()
		preparedRequest, err := db.Prepare("SELECT MAX(cid) FROM Comments;")
		checkErr(err)
		var maxCid int
		preparedRequest.QueryRow().Scan(&maxCid)
		commentData.Cid = maxCid + 1

		_, err = db.Exec("INSERT INTO Comments (cid, content, like, createdDate, uid, pid) VALUES (?, ?, ?, ?, ?, ?);", commentData.Cid, commentData.Content, commentData.Like, commentData.CreatedDate, commentData.Uid, commentData.Pid)
		checkErr(err)
		commentData.Content = r.FormValue("content")
	}


	// refresh
	publicationData := makePublicationWithId(commentData.Pid, w, r, "addCommentBox")

	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))

	err = tpl.Execute(w, publicationData)
	checkErr(err)
}