package publications

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type publicationTemplateData struct {
	IdPublication int
	Title         string
	Description   string
	ImageLink     string
	UpvoteNumber  int
	CreatedDate   string
	UsernameId    int

	Username      string
	CommentNumber int
	IsThereImage  bool
	Tags          template.HTML
}

/*
Take the id of a publication to give a 70% wide and 150px tall card of the publication
*/
func MakePublicationHomePageTemplate(idPublication string) template.HTML {
	publicationTemplate := publicationTemplateData{}

	db, err := sql.Open("sqlite3", "./database.db")
	checkErr(err)
	defer db.Close()

	// Get the publication from the db
	preparedRequest, err := db.Prepare("SELECT * FROM Publications WHERE pid = ?;")
	checkErr(err)
	row := preparedRequest.QueryRow(idPublication)
	row.Scan(&publicationTemplate.IdPublication, &publicationTemplate.Title, &publicationTemplate.Description, &publicationTemplate.ImageLink,
		&publicationTemplate.UpvoteNumber, &publicationTemplate.CreatedDate, &publicationTemplate.UsernameId)

	// isThereImage
	if publicationTemplate.ImageLink != "" {
		publicationTemplate.IsThereImage = true
	} else {
		publicationTemplate.IsThereImage = false
	}

	// get Username
	preparedRequest, err = db.Prepare("SELECT username FROM Users WHERE uid = ?;")
	checkErr(err)
	preparedRequest.QueryRow(publicationTemplate.UsernameId).Scan(&publicationTemplate.Username)

	// get number of comment
	preparedRequest, err = db.Prepare("SELECT COUNT(*) FROM Comments WHERE pid = ?;")
	checkErr(err)
	preparedRequest.QueryRow(publicationTemplate.IdPublication).Scan(&publicationTemplate.CommentNumber)

	// get tags
	preparedRequest, err = db.Prepare("SELECT name FROM Tags WHERE pid = ?")
	checkErr(err)
	rows, err := preparedRequest.Query(publicationTemplate.IdPublication)
	checkErr(err)
	var tagArray []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		tagArray = append(tagArray, tag)
	}
	publicationTemplate.Tags = makeTags(tagArray)

	tpl := new(bytes.Buffer)
	tplRaw := template.Must(template.ParseFiles("templates/publicationTemplate.html"))
	err = tplRaw.Execute(tpl, publicationTemplate)
	checkErr(err)
	tplString := tpl.String()
	return template.HTML(tplString)
}

func makeTags(tags []string) template.HTML {
	finalString := ""
	for _, tag := range tags {
		finalString += "<div class=\"publicationTag\" style=\"background-color: "

		switch tag {
		case "Gaming":
			finalString += "#0033cc\">" // blue
			break
		case "Lifestyle":
			finalString += "#ff3399\">" // pink
			break
		case "Space":
			finalString += "#000066\">" // dark blue
			break
		case "Art":
			finalString += "#ff3300\">" // red
			break
		case "Nature":
			finalString += "#009933\">" // green
			break
		default:
			finalString += "#000000\">" // black
		}
		finalString += tag + "</div>"

	}

	return template.HTML(finalString)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
