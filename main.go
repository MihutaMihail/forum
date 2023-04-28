package main

import (
	"fmt"
	"forum/code/publications"
	"html/template"
	"log"
	"net/http"
)

type indexPageData struct {
	Publication01 template.HTML
	Publication02 template.HTML
}

func main() {

	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))

	indexData := indexPageData{
		Publication01: publications.MakePublicationHomePageTemplate("title", "description blabla", "http://www.snut.fr/wp-content/uploads/2015/08/image-de-paysage-2.jpg"),
		Publication02: publications.MakePublicationHomePageTemplate("title2", "c'est une deuxième description", "https://bonjournature.fr/wp-content/uploads/image-automne-gratuite-1.jpg"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePage := template.Must(template.ParseFiles("templates/index.html"))

		homePage.Execute(w, indexData)
	})

	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
