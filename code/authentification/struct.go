package authentification

import "text/template"

type Server struct {
	Tpl *template.Template
}

type DataAuthentification struct {
	Id                 int
	IsCookie           bool
	ErreurCookie string
	SaisieMail         []string
	ErreurVoid         string
	ErreurMail         string
	Mail               string
	PasswordBeforeHash string
	PasswordAfterHash  string
	Password           string
	Username           string
}
