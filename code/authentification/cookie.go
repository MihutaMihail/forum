package authentification

import (
	"fmt"
	"log"
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

/*
Return 0 if not connected
*/
func GetSessionUid(w http.ResponseWriter, r *http.Request) int {
	c, err := r.Cookie("session")
	if err == nil {
		uid, err := strconv.Atoi(c.Value)
		if err != nil {
			log.Panic(err)
		}
		return uid
	} else {
		return 0
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
