package publications

import (
	"bytes"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

/*
Take the id of a publication to give a 70% wide and 150px tall card of the publication
*/
func MakePublicationHomePageTemplate(idPublication int, w http.ResponseWriter, r *http.Request) template.HTML {

	publicationData := makePublicationWithId(idPublication, w, r)

	tpl := new(bytes.Buffer)
	tplRaw := template.Must(template.ParseFiles("templates/publicationTemplate.html"))
	err := tplRaw.Execute(tpl, publicationData)
	checkErr(err)
	tplString := tpl.String()
	return template.HTML(tplString)
}

func MakeTags(tags []string) template.HTML {
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
