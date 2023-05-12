package publications

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"strconv"
)

type PublicationData struct {
	Pid          int
	Title        string
	Content      string
	ImageLink    string
	UpvoteNumber int
	CreatedDate  string
	Uid          string

	IsThereImage       bool
	Username           string
	CommentNumber      int
	Tags               template.HTML
	TagsString         []string
	Comments           []CommentData
	SortedByPertinance bool

	UpvoteClass      string
	DownvoteClass    string
	CreateCommentBox template.HTML
}
type CommentData struct {
	Cid         int
	Content     string
	Like        int
	CreatedDate string
	Uid         int
	Pid         int

	Username string

	UpvoteClass   string
	DownvoteClass string
}

const commentBoxTemplateFirst string = "<div class=\"commentBox\"><form method=\"POST\" action=\"/sendComment?pid="
const commentBoxTemplateSecond string = "\">	<textarea class=\"commentTyping\" name=\"content\"></textarea><input class=\"commentSend\" type=\"submit\"></form></div>"

/*
Accept any args :
  - addCommentBox : add the comment box to write a comment
*/
func makePublicationWithId(idInt int, args ...string) *PublicationData {
	// open the db
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	publicationData := PublicationData{}

	for _, arg := range args {
		switch arg {
		case "addCommentBox": 
			publicationData.CreateCommentBox = template.HTML(commentBoxTemplateFirst + strconv.Itoa(idInt) + commentBoxTemplateSecond)
			break;
		default:
			fmt.Println("Warning : invalid arg at call to makePublicationWithId")
		}
	}

	// Get the publication from the db
	preparedRequest, err := db.Prepare("SELECT * FROM Publications WHERE pid = ?;")
	checkErr(err)
	row := preparedRequest.QueryRow(idInt)
	row.Scan(&publicationData.Pid, &publicationData.Title, &publicationData.Content, &publicationData.ImageLink, &publicationData.UpvoteNumber, &publicationData.CreatedDate, &publicationData.Uid)

	// isThereImage
	if publicationData.ImageLink != "" {
		publicationData.IsThereImage = true
	} else {
		publicationData.IsThereImage = false
	}

	// get username
	preparedRequest, err = db.Prepare("SELECT username FROM Users WHERE uid = ?;")
	checkErr(err)
	preparedRequest.QueryRow(publicationData.Uid).Scan(&publicationData.Username)

	// get number of comment
	preparedRequest, err = db.Prepare("SELECT COUNT(*) FROM Comments WHERE pid = ?;")
	checkErr(err)
	preparedRequest.QueryRow(publicationData.Pid).Scan(&publicationData.CommentNumber)

	// get tags
	preparedRequest, err = db.Prepare("SELECT name FROM Tags WHERE pid = ?;")
	checkErr(err)
	rows, err := preparedRequest.Query(publicationData.Pid)
	defer rows.Close()
	checkErr(err)
	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		tagArray = append(tagArray, tag)
	}
	publicationData.Tags = MakeTags(tagArray)

	//liked or not by session user
	uid := 1 // getSessionUid                 // TODO
	preparedRequest, err = db.Prepare("SELECT isLike FROM Likes WHERE uid = ? AND (pid = ? AND pid != 0);")
	checkErr(err)
	rows, err = preparedRequest.Query(uid, publicationData.Pid)
	defer rows.Close()
	for rows.Next() { // if there is a like, it will do one loop, else it will pass
		var isLike int
		err = rows.Scan(&isLike)
		checkErr(err)
		if isLike == 1 { // upvote or downvote
			publicationData.UpvoteClass = "clickedVote"
		} else {
			publicationData.DownvoteClass = "clickedVote"
		}
	}

	publicationData.Comments = makeComments(publicationData.Pid)

	// fmt.Println(publicationData.Comments[0].Content)
	return &publicationData
}

func makeComments(Pid int) []CommentData {
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()
	var finalArray []CommentData

	// get all comments
	preparedRequest, err := db.Prepare("SELECT * FROM Comments WHERE pid = ?;")
	checkErr(err)
	rows, err := preparedRequest.Query(Pid)
	defer rows.Close()
	// for each results, get the comment data
	for rows.Next() {
		var comment CommentData
		err := rows.Scan(&comment.Cid, &comment.Content, &comment.Like, &comment.CreatedDate, &comment.Uid, &comment.Pid)
		checkErr(err)

		// get username of comment
		preparedRequest, err = db.Prepare("SELECT username FROM Users WHERE uid = ?;")
		checkErr(err)
		preparedRequest.QueryRow(comment.Uid).Scan(&comment.Username)

		//liked or not by session user
		uid := 1 // getSessionUid                 // TODO
		preparedRequest, err = db.Prepare("SELECT isLike FROM Likes WHERE uid = ? AND (cid = ? AND cid != 0);")
		checkErr(err)
		rowsLike, err := preparedRequest.Query(uid, comment.Cid)
		defer rowsLike.Close()
		for rowsLike.Next() { // if there is a like, it will do one loop, else it will pass
			var isLike int
			err = rowsLike.Scan(&isLike)
			checkErr(err)
			if isLike == 1 { // upvote or downvote
				comment.UpvoteClass = "clickedVote"
			} else {
				comment.DownvoteClass = "clickedVote"
			}
		}

		finalArray = append(finalArray, comment)
	}

	return finalArray
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
