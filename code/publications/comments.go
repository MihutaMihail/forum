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

	indexData := indexPageData{}
	indexData.Main = parsePublicationPage(w, r, true, pid)
	indexData.Header = MakeHeaderTemplate(w, r, false)

	tpl := template.Must(template.ParseFiles("templates/publicationListTemplate.html"))

	err = tpl.Execute(w, indexData)
	checkErr(err)
}

/*
Handle to add a comment in the final click
*/
func AddAComment(w http.ResponseWriter, r *http.Request) {
	var err error
	commentData := CommentData{}

	cookie := authentification.GetSessionUid(w, r)

	commentData.Pid, err = strconv.Atoi(r.URL.Query().Get("pid"))

	if cookie != 0 {
		commentData.Uid = cookie
		
		checkErr(err)
		commentData.Like = 0

		timeNow := time.Now()
		commentData.CreatedDate = timeNow.Format("02/01/2006")

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
	refreshPublicationPage(w, r, commentData.Pid)
}


func CommentSortPertinance(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	checkErr(err)
	
	deleteCommentSortedByDateCookie(w, r)
	http.Redirect(w, r, "/publication?pid=" + strconv.Itoa(pid), http.StatusSeeOther)
}

func CommentSortDate(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	checkErr(err)

	setCommentSortedByDateCookie(w, r)
	http.Redirect(w, r, "/publication?pid=" + strconv.Itoa(pid), http.StatusSeeOther)
}

func setCommentSortedByDateCookie(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(6 * time.Hour)
	cookie, err := r.Cookie("sortingByDate")
	if err != nil {
		cookie = &http.Cookie{Name: "sortingByDate", Value: "true", Expires: expiration}
		http.SetCookie(w, cookie)
	}
	// I'm hungry now
}

func deleteCommentSortedByDateCookie(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(-time.Hour) // expiration in the past, it delete
	cookie := &http.Cookie{Name: "sortingByDate", Value: "false", Expires: expiration}

	http.SetCookie(w, cookie)
}