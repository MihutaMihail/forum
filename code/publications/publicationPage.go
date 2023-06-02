package publications

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
) 

func HandlePublication(w http.ResponseWriter, r *http.Request) {
	indexData := indexPageData{}
	indexData.Main = parsePublicationPage(w, r, false)
	indexData.Header = MakeHeaderTemplate(w, r, false)

	tpl := template.Must(template.ParseFiles("templates/publicationListTemplate.html"))

	err := tpl.Execute(w, indexData)
	checkErr(err)
}

func parsePublicationPage(w http.ResponseWriter, r *http.Request, commentBox bool, pid ...int) template.HTML{
	// get id and the corresponding data by url or signature
	var id int
	var err error
 	var publicationData *PublicationData

	r.ParseForm()
	if (len(pid) != 0) {
		id = pid[0]
	} else {
		id, err = strconv.Atoi(r.URL.Query().Get("pid"))
		checkErr(err)
	}

	// if the publication is made with a comment box
	if commentBox {
		publicationData = makePublicationWithId(id, w, r, "addCommentBox")
	} else {
		publicationData = makePublicationWithId(id, w, r)
	}

	tpl := new(bytes.Buffer)
	tplRaw := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	err = tplRaw.Execute(tpl, publicationData)
	checkErr(err)
	tplString := tpl.String()
	return template.HTML(tplString)
}
