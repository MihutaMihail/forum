package publications

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type PublicationTemplateData struct {
	IdPublication string
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
	And it's at this moment that Fafa realized. There is no overloading (or even default variables) in go.

So, PLEASE, pass an empty string if the post don't have an image, and sorry

Will only need the id to access the database in the future
*/
func MakePublicationHomePageTemplate(idPublication string) template.HTML {
	publicationTemplate := PublicationTemplateData{}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the publication from the db
	preparedRequest, err := db.Prepare("SELECT * FROM Publications WHERE pid = ?")
	if err != nil {
		log.Fatal(err)
	}
	rows := preparedRequest.QueryRow(idPublication)
	rows.Scan(&publicationTemplate.IdPublication, &publicationTemplate.Title, &publicationTemplate.Description, &publicationTemplate.ImageLink,
		&publicationTemplate.UpvoteNumber, &publicationTemplate.CreatedDate, &publicationTemplate.UsernameId)
	fmt.Println(publicationTemplate.ImageLink)

	if publicationTemplate.ImageLink != "" {
		publicationTemplate.IsThereImage = true
	} else {
		publicationTemplate.IsThereImage = false
	}
	publicationTemplate.CommentNumber = 10 // TEMP

	tpl := new(bytes.Buffer)

	tplRaw := template.Must(template.ParseFiles("templates/publicationTemplate.html"))

	err = tplRaw.Execute(tpl, publicationTemplate)

	if err != nil {
		log.Fatal(err)
	}

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
