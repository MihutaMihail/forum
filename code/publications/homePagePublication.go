package publications

import (
	"bytes"
	"html/template"
	"log"
)

type publicationTemplateData struct {
	Title         string
	Description   string
	UpvoteNumber  int
	CommentNumber int
	ImageLink     string
	IsThereImage  bool
}

func MakePublicationHomePageTemplate(title string, description string, imageLink string) template.HTML {
	publicationTemplate := publicationTemplateData{
		Title:        title,
		Description:  description,
		UpvoteNumber: 0,
		ImageLink:    imageLink,
		IsThereImage: true,
	}

	tpl := new(bytes.Buffer)

	tplRaw := template.Must(template.ParseFiles("templates/publicationTemplate.html"))

	err := tplRaw.Execute(tpl, publicationTemplate)
	if err != nil {
		log.Fatal(err)
	}

	tplString := tpl.String()
	t := template.HTML(tplString)
	return t

}
