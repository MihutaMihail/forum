package authentification

import (
	"net/http"
	"text/template"
)

var server Server

func init() {
	server.Tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	Open()
	server.Tpl.ExecuteTemplate(w, "login.html", Data)
}

func Register(w http.ResponseWriter, r *http.Request) {
	Open()
	server.Tpl.ExecuteTemplate(w, "register.html", Data)
}