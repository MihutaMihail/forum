package publications

import (
	"database/sql"
	"html/template"
	"log"
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
	Comments           []CommentData
	SortedByPertinance bool
}
type CommentData struct {
	Cid         int
	Content     string
	Like        int
	CreatedDate string
	Uid         int
	Pid         int

	Username string
}

func makePublicationWithId(idInt int) *PublicationData{
	// open the db
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	publicationData := PublicationData{}

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
	checkErr(err)
	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		tagArray = append(tagArray, tag)
	}
	publicationData.Tags = MakeTags(tagArray)

	publicationData.Comments = makeComments(publicationData.Pid)

	// fmt.Println(publicationData.Comments[0].Content)
	return &publicationData
}

func makeComments(Pid int) []CommentData{
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()
	var finalArray []CommentData

	// get all comments
	preparedRequest, err := db.Prepare("SELECT * FROM Comments WHERE pid = ?;") 
	checkErr(err)
	rows, err := preparedRequest.Query(Pid)
	// for each results, get the comment data
	for rows.Next() {
		var comment CommentData
		err := rows.Scan(&comment.Cid, &comment.Content, &comment.Like, &comment.CreatedDate, &comment.Uid, &comment.Pid)
		checkErr(err)

		// get username of comment
		preparedRequest, err = db.Prepare("SELECT username FROM Users WHERE uid = ?;")
		checkErr(err)
		preparedRequest.QueryRow(comment.Uid).Scan(&comment.Username)

		
		finalArray = append(finalArray, comment)
	}

	return finalArray
}


func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}