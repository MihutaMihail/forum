package publications

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
)

type publicationData struct {
	pid          int
	title        string
	content      string
	imageLink    string
	upvoteNumber int
	createdDate  string
	uid          string

	isThereImage       bool
	username           string
	commentNumber      int
	tags               template.HTML
	comments           []commentData
	sortedByPertinance bool
}
type commentData struct {
	cid         int
	content     string
	like        int
	createdDate string
	uid         int
	pid         int

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
	publicationData.pid, err = strconv.Atoi(r.FormValue("idPublication"))
	checkErr(err)

	// Get the publication from the db
	preparedRequest, err := db.Prepare("SELECT * FROM Publications WHERE pid = ?;")
	checkErr(err)
	row := preparedRequest.QueryRow(publicationData.pid)
	row.Scan(&publicationData.pid, &publicationData.title, &publicationData.content, &publicationData.imageLink, &publicationData.upvoteNumber, &publicationData.createdDate, &publicationData.uid)



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
	
	tpl.Execute(w, publicationData)
}
