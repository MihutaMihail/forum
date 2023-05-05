package publications

import (
	"html/template"
	"net/http"
	"strconv"
)

func HandlePublication(w http.ResponseWriter, r *http.Request) /*template.HTML*/ {
	if r.Method != "POST" {
		http.Error(w, "Expected POST, found "+r.Method, http.StatusBadRequest) // 400
		return /*template.HTML("err 400: bad Request")*/
	}

	// get id and the corresponding data
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("idPublication"))
	checkErr(err)
	publicationData := makePublicationWithId(id)
	


		//
	// create the template and returns it. For when the handler will change to add the website interface
	//
	// tpl := new(bytes.Buffer)
	// tplRaw := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	// err = tplRaw.Execute(tpl, publicationData)
	// checkErr(err)
	// tplString := tpl.String()
	// return template.HTML(tplString)
	
	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))
	
	err = tpl.Execute(w, publicationData)
	checkErr(err)
}
