package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePage := template.Must(template.ParseFiles("templates/index.html"))

		homePage.Execute(w, "o")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
