package main

import (
	"database/sql"
	"fmt"
	"forum/code/publications"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type indexPageData struct {
	Publication01 template.HTML
	Publication02 template.HTML
	Publication03 template.HTML
}

func main() {
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))

	// test()

	indexData := indexPageData{
		Publication01: publications.MakePublicationHomePageTemplate(1),
		Publication02: publications.MakePublicationHomePageTemplate(2),
		Publication03: publications.MakePublicationHomePageTemplate(3),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePage := template.Must(template.ParseFiles("templates/index.html"))

		err := homePage.Execute(w, indexData)
		if err!=nil{log.Panic(err)}
	})
	http.HandleFunc("/publication", publications.HandlePublication)

	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// here you can do dynamic tests
func test() {

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		log.Fatal(err)
	}
	test := ""

	for rows.Next() {
		err := rows.Scan(&test)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(test)
	}
}
