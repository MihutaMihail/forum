package authentification

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

var Data DataAuthentification

func Open() {
	db, _ = sql.Open("sqlite3", "./database.db")
}

func Print(w http.ResponseWriter, r *http.Request) {
	row, _ := db.Query("SELECT uid, email, password, username FROM Users")

	for row.Next() {
		row.Scan(&Data.Id, &Data.Mail, &Data.Password, &Data.Username)
		fmt.Println(strconv.Itoa(Data.Id) + " Mail : " + Data.Mail + ", Password : " + Data.Password + ", Username : " + Data.Username)
	}
	row.Close()
	println()

	http.Redirect(w, r, "http://localhost:8080/register", http.StatusSeeOther)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	statement, _ := db.Prepare("DELETE FROM Users")
	statement.Exec()
	println()
	println("Base de donnée reset")
	println()
	http.Redirect(w, r, "http://localhost:8080/register", http.StatusSeeOther)
}

func AddUsers(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("mail")) == 0 || len(r.FormValue("password")) == 0 || len(r.FormValue("username")) == 0 {
		Data.ErreurVoid = "Un des champs est vide"
		http.Redirect(w, r, "http://localhost:8080/register", http.StatusSeeOther)
	} else {
	existingUser, err := db.Query("SELECT uid FROM Users WHERE email = ? OR username = ?", r.FormValue("mail"), r.FormValue("username"))
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	if !existingUser.Next() {
		add, err := db.Prepare("INSERT INTO Users (email, password, username, admin) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = add.Exec(r.FormValue("mail"), r.FormValue("password"), r.FormValue("username"), 0)
		if err != nil {
			log.Fatal(err)
		}
		existingUser.Close()
		Data.ErreurMail = ""
		Data.ErreurVoid = ""
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
	} else {
		existingUser.Close()
		Data.ErreurMail = "Le mail ou username est déja pris"
		http.Redirect(w, r, "http://localhost:8080/register", http.StatusSeeOther)
	}
}
}

func CheckUsers(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("mail")) == 0 || len(r.FormValue("password")) == 0 || len(r.FormValue("username")) == 0 {
		Data.ErreurVoid = "Un des champs est vide"
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
	} else {
			existingUser, err := db.Query("SELECT uid FROM Users WHERE email = ? AND password = ? AND username = ?", r.FormValue("mail"), r.FormValue("password"), r.FormValue("username"))
		if err != nil {
			panic(err)
		}

		if existingUser.Next() {
			println("LOGIN MARCHE")
			existingUser.Scan(&Data.Id)
			SetSessionUid(w,r,Data.Id)
			existingUser.Close()
			http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
		} else {
			existingUser.Close()
			Data.ErreurMail = "Le mail, le mot de passe ou l'username ne correspond pas"
			println("LOGIN MARCHE PAS")
			http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
		}
	}
		

		
	}
