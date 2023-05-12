package authentification

import (
	"net/http"
)

func RegisterGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "http://localhost:8082/register", http.StatusSeeOther)
		return
	}

	AddUsers(w, r)
}