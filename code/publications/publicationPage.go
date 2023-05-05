package publications

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
)

type publicationData struct {
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
	Comments           []commentData
	SortedByPertinance bool
}
type commentData struct {
	Cid         int
	Content     string
	Like        int
	CreatedDate string
	Uid         int
	Pid         int

	username string
}

func HandlePublication(w http.ResponseWriter, r *http.Request) /*template.HTML*/ {
	if r.Method != "POST" {
		http.Error(w, "Expected POST, found "+r.Method, http.StatusBadRequest) // 400
		return /*template.HTML("err 400: bad Request")*/
	}
	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	publicationData := publicationData{}

	r.ParseForm()
	publicationData.Pid, err = strconv.Atoi(r.FormValue("idPublication"))
	checkErr(err)

	// Get the publication from the db
	preparedRequest, err := db.Prepare("SELECT * FROM Publications WHERE pid = ?;")
	checkErr(err)
	row := preparedRequest.QueryRow(publicationData.Pid)
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
	preparedRequest, err = db.Prepare("SELECT name FROM Tags WHERE pid = ?")
	checkErr(err)
	rows, err := preparedRequest.Query(publicationData.Pid)
	checkErr(err)
	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		tagArray = append(tagArray, tag)
	}
	publicationData.Tags = makeTags(tagArray)


		//
	// create the template and returns it. For when the handler will change to add the website interface
	//
	// tpl := new(bytes.Buffer)
	// tplRaw := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	// err = tplRaw.Execute(tpl, publicationData)
	// checkErr(err)
	// tplString := tpl.String()
	// return template.HTML(tplString)
	
	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	
	err = tpl.Execute(w, publicationData)
	checkErr(err)
}
