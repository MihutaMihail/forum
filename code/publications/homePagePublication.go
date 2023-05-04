package publications

import (
	"bytes"
	"html/template"
	"log"
)

type PublicationTemplateData struct {
	IdPublication string
	Username      string
	Title         string
	Description   string
	Tags          template.HTML
	UpvoteNumber  int
	CommentNumber int
	ImageLink     string
	IsThereImage  bool
}

/*
	And it's at this moment that Fafa realized. There is no overloading (or even default variables) in go.

So, PLEASE, pass an empty string if the post don't have an image, and sorry

Will only need the id to access the database in the future
*/
func MakePublicationHomePageTemplate(idPublication string, username string, title string, description string, imageLink string, tags []string, upvoteNumber int, commentNumber int) template.HTML {
	publicationTemplate := PublicationTemplateData{
		IdPublication: idPublication,
		Username:      username,
		Title:         title,
		Description:   description,
		Tags:          makeTags(tags),
		UpvoteNumber:  upvoteNumber,
		ImageLink:     imageLink,
		IsThereImage:  imageLink != "",
		CommentNumber: commentNumber,
	}

	tpl := new(bytes.Buffer)

	tplRaw := template.Must(template.ParseFiles("templates/publicationTemplate.html"))

	err := tplRaw.Execute(tpl, publicationTemplate)
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