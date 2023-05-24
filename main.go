package main

import (
	"fmt"
	"forum/code/authentification"
	"forum/code/publications"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	// TEMPORARY
	http.Handle("/assets/uploads/", http.StripPrefix("/assets/uploads/", http.FileServer(http.Dir("assets/uploads"))))
	// TEMPORARY

	http.HandleFunc("/", publications.HandleAllPosts)
	http.HandleFunc("/publication", publications.HandlePublication)
	http.HandleFunc("/likes", publications.HandleLikes)
	http.HandleFunc("/addCommentBox", publications.MakeCommentBox)
	http.HandleFunc("/sendComment", publications.AddAComment)
	http.HandleFunc("/publicationForm", publications.HandleFormPost)
	http.HandleFunc("/checkpublicationForm", publications.CheckHandleFormPost)
	http.HandleFunc("/publicationSubmitForm", publications.HandleSubmitForm)
	http.HandleFunc("/publicationDelete", publications.HandleDeletePost)

	http.HandleFunc("/login", authentification.Login)
	http.HandleFunc("/loginGet", authentification.LoginGet)
	http.HandleFunc("/register", authentification.Register)
	http.HandleFunc("/print", authentification.Print)
	http.HandleFunc("/registerGet", authentification.RegisterGet)

	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// here you can do dynamic tests
/*func test() {

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
}*/
