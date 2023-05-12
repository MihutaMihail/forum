package publications

import (
	"html/template"
	"net/http"
	"strconv"
)

/*
Handle for the first click to add a comment ; call the same page with a commentBox
*/
func MakeCommentBox(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	checkErr(err)
	publicationData := makePublicationWithId(pid, "addCommentBox")

	tpl := template.Must(template.ParseFiles("templates/publicationPageTemplate.html"))

	err = tpl.Execute(w, publicationData)
	checkErr(err)
}

/*
Handle to add a comment in the final click
*/
func AddAComment(w http.ResponseWriter, r *http.Request) {

}
