package main

import (
	"database/sql"
	"fmt"
	"forum/code/publications"
	"forum/code/testcrud"
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

	http.HandleFunc("/", testcrud.HandleAllPosts)
	

	 //
    // TEST CRUD
    //
    http.Handle("/testcrud/uploads/", http.StripPrefix("/testcrud/uploads/", http.FileServer(http.Dir("testcrud/uploads"))))

    http.HandleFunc("/testcrud", testcrud.HandleIndex)
    http.HandleFunc("/testFormPost", testcrud.HandleFormPost)
    http.HandleFunc("/testSubmitForm", testcrud.HandleSubmitForm)
    http.HandleFunc("/testPost", testcrud.HandlePost)
    http.HandleFunc("/testDeletePost", testcrud.HandleDeletePost)

    // TEST CRUD
	http.HandleFunc("/publication", publications.HandlePublication)
	http.HandleFunc("/likes", publications.HandleLikes)

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
