package publications

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
) 

func HandlePublication(w http.ResponseWriter, r *http.Request) {
	indexData := indexPageData{}
	indexData.Main = parsePublicationPage(w, r)
	indexData.Header = MakeHeaderTemplate(w, r)

	tpl := template.Must(template.ParseFiles("templates/publicationListTemplate.html"))

	err := tpl.Execute(w, indexData)
	checkErr(err)
}

func parsePublicationPage(w http.ResponseWriter, r *http.Request) template.HTML{
	// get id and the corresponding data
	r.ParseForm()
	id, err := strconv.Atoi(r.URL.Query().Get("pid"))
	checkErr(err)
	publicationData := makePublicationWithId(id, w, r)

	tpl := new(bytes.Buffer)
	tplRaw := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	err = tplRaw.Execute(tpl, publicationData)
	checkErr(err)
	tplString := tpl.String()
	return template.HTML(tplString)
}
