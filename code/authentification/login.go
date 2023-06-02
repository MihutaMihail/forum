package authentification

import (
	"net/http"
	"log"
)

func LoginGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
	}

	if CheckSessionUid(w,r) != nil {
		CheckUsers(w, r)
	} else {
		Data.ErreurCookie = "Vous avez déjà votre session ouverte"
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)

	}	
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := DeleteSessionCookie(w, r) 
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
