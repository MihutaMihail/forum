package authentification

import "net/http"

func LoginGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
	}

	if CheckSessionUid(w,r) {
		CheckUsers(w, r)
	} else {
		Data.ErreurCookie = "Vous avez déjà votre session ouverte"
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)

	}	
}