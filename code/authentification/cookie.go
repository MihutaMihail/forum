package authentification

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func SetSessionUid(w http.ResponseWriter, r *http.Request, id int) {
	e := time.Now().Add(6 * time.Hour)
	uid := strconv.Itoa(id)
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{Name: "session", Value: uid, Expires: e}
		http.SetCookie(w, c)
	} else {
		fmt.Println(w, "Le cookie est déja la !")
	}
}

func GetSessionUid(w http.ResponseWriter, r *http.Request) *http.Cookie {
	c, err := r.Cookie("session")
	if err == nil {
		return c
	} else {
		return nil
	}
}

func CheckSessionUid(w http.ResponseWriter, r *http.Request) error {
	_, err := r.Cookie("session")
	if err != nil {
		return err
	} else {
		return nil
	}
}