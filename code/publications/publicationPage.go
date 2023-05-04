package publications

import (
	"html/template"
	"net/http"
)

func HandlePublication(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Expected POST, found "+r.Method, http.StatusBadRequest) // 400
		return
	}
	r.ParseForm()

	idPublication := r.FormValue("idPublication")

	http.Error(w, "501 ; WIP", http.StatusNotImplemented) //501

	testString := template.Must(template.New("name").Parse(idPublication + ", it works !\n"))
	testString.Execute(w, "o")

	// TODO ; make the publication page by getting the data from database + a template needed

	return
}
