package publications

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type LikeData struct {
	Uid    int
	Pid    int // if 0 ; it's a comment
	Cid    int // if 0 ; it's a publication
	IsLike bool
}

/*
request url be like ; /likes?id=25&isComment=true&isLike=false
*/
func HandleLikes(w http.ResponseWriter, r *http.Request) {
	likeData := LikeData{Pid: 0, Cid: 0}

	//uid := authentification.GetSessionUid(w,r) // TODO
	
	likeData.Uid = 1 //uid
	var err error

	// get pid or cid
	if r.URL.Query().Get("isComment") == "true" {
		likeData.Cid, err = strconv.Atoi(r.URL.Query().Get("id"))
		checkErr(err)
	} else {
		likeData.Pid, err = strconv.Atoi(r.URL.Query().Get("id"))
		checkErr(err)
	}

	// get isLike
	if r.URL.Query().Get("isLike") == "true" {
		likeData.IsLike = true
	} else {
		likeData.IsLike = false
	}

	// open the db
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// now we check if there is already a like from user
	preparedRequest, err := db.Prepare("SELECT COUNT(*) FROM Likes WHERE uid = ? AND ((pid = ? AND pid != 0) OR (cid = ? AND cid != 0))")
	checkErr(err)
	var tempNbLikes int
	preparedRequest.QueryRow(likeData.Uid, likeData.Pid, likeData.Cid).Scan(&tempNbLikes)

	if tempNbLikes > 1 { // would be not good
		fmt.Println("WARNING : more than 1 like/dislike from user " + strconv.Itoa(likeData.Uid) + " at the publication/comment " + strconv.Itoa(likeData.Pid) + "/" + strconv.Itoa(likeData.Cid))
	}
	if tempNbLikes == 1 { // already has a like/dislike
		alreadyALike(likeData, db)
	}
	if tempNbLikes == 0 { // doesn't have like/dislike, easy, just need to add
		addLikeOrDislike(likeData, db)
	}

	// if we're not on a publi, get the publi before refreshing it
	if likeData.Pid == 0 {
		likeData.Pid = getPidFromCid(likeData, db)
	}
	// recreate the publication (it refresh)
	publicationData := makePublicationWithId(likeData.Pid)
	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	err = tpl.Execute(w, publicationData)
	checkErr(err)
}

func alreadyALike(likeData LikeData, db *sql.DB) {
	preparedRequest, err := db.Prepare("SELECT isLike FROM Likes WHERE uid = ? AND ((pid = ? AND pid != 0) OR (cid = ? AND cid != 0))")
	checkErr(err)
	var tempOldInteractionInt int
	preparedRequest.QueryRow(likeData.Uid, likeData.Pid, likeData.Cid).Scan(&tempOldInteractionInt)
	var oldInteraction bool
	if tempOldInteractionInt == 1 {
		oldInteraction = true
	} else {
		oldInteraction = false
	}

	// if user click a second time on the same arreow, else, he clicked on the other one
	if likeData.IsLike == oldInteraction { // then we should remove the like/dislike
		_, err = db.Exec("DELETE FROM Likes WHERE uid = ? AND ((pid = ? AND pid != 0) OR (cid = ? AND cid != 0))", likeData.Uid, likeData.Pid, likeData.Cid)
		checkErr(err)

		updateLikeCounter(likeData, db)
	} else { // then we should switch from like to dislike, or other way around
		_, err = db.Exec("UPDATE Likes SET isLike = ? WHERE uid = ? AND ((pid = ? AND pid != 0) OR (cid = ? AND cid != 0))", !oldInteraction, likeData.Uid, likeData.Pid, likeData.Cid)
		checkErr(err)

		updateLikeCounter(likeData, db)
	}
}

func addLikeOrDislike(likeData LikeData, db *sql.DB) {
	_, err := db.Exec("INSERT INTO Likes (uid, pid, cid, isLike) VALUES (?, ?, ?, ?)", likeData.Uid, likeData.Pid, likeData.Cid, likeData.IsLike)
	checkErr(err)

	updateLikeCounter(likeData, db)
}

/*
 change the likeCounter on the publication or comment by the amount in case of like, -amount in case of dislike
*/
// let's just count
func updateLikeCounter(likeData LikeData, db *sql.DB) {
	var finalResult int
	var rows *sql.Rows
	var err error

	if likeData.Cid == 0 {
		preparedRequest, err := db.Prepare("SELECT isLike FROM Likes WHERE pid = ?")
		checkErr(err)
		rows, err = preparedRequest.Query(likeData.Pid)
	} else {
		preparedRequest, err := db.Prepare("SELECT isLike FROM Likes WHERE cid = ?")
		checkErr(err)
		rows, err = preparedRequest.Query(likeData.Cid)
	}
	for rows.Next() {
		var isLikeInt int
		err = rows.Scan(&isLikeInt)
		checkErr(err)
		if isLikeInt == 1 {
			finalResult++
		} else {
			finalResult--
		}
	}

	if likeData.Cid == 0 {
		_, err = db.Exec("UPDATE Publications SET like = ? WHERE pid = ?", finalResult, likeData.Pid)
	} else {
		_, err = db.Exec("UPDATE Comments SET like = ? WHERE cid = ?", finalResult, likeData.Cid)
	}
	checkErr(err)
}

/*
self-explanatory
*/
func getPidFromCid(likeData LikeData, db *sql.DB) int {
	preparedRequest, err := db.Prepare("SELECT pid FROM Comments WHERE cid = ?")
	checkErr(err)
	var pid int
	preparedRequest.QueryRow(likeData.Cid).Scan(&pid)
	return pid
}
