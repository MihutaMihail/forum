package publications

import (
	"html/template"
	"net/http"
	"strconv"
)

func HandlePublication(w http.ResponseWriter, r *http.Request) {

	// get id and the corresponding data
	r.ParseForm()
	id, err := strconv.Atoi(r.URL.Query().Get("pid"))
	checkErr(err)
	publicationData := makePublicationWithId(id, w, r)


	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))

	err = tpl.Execute(w, publicationData)
	checkErr(err)
}
