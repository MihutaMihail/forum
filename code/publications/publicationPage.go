package publications

import (
	"html/template"
	"net/http"
	"strconv"
)

func HandlePublication(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Expected POST, found "+r.Method, http.StatusBadRequest) // 400
		return                                                                 /*template.HTML("err 400: bad Request")*/
	}

	// get id and the corresponding data
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("idPublication"))
	checkErr(err)
	publicationData := makePublicationWithId(id, w, r)


	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))

	err = tpl.Execute(w, publicationData)
	checkErr(err)
}
