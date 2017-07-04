package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

var tpl *template.Template
var db *sql.DB
var err error
var store = sessions.NewCookieStore([]byte("Lazy dog jumped over fire"))

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	db, err = sql.Open("mysql", "<DSN>")
	check(err)
	defer db.Close()
	err = db.Ping()
	check(err)

	// setup gomniauth
	gomniauth.SetSecurityKey("The lazy dog jumped over fire")
	gomniauth.WithProviders(
		google.New("607540849947-8mrgd1unre20gkfqtjfo72uqc1b3fec0.apps.googleusercontent.com",
			"0q0YcklzYSOvFvKsCGhOtI1J",
			"http://localhost:8080/callback"),
	)

	http.HandleFunc("/", index)
	http.HandleFunc("/signupform", signupform)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/completeform", completeform)
	http.HandleFunc("/complete", complete)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logingoogle", logingoogle)
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/editform", editform)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/forgot", forgot)
	http.HandleFunc("/reset", reset)
	http.HandleFunc("/resetpasswordform", resetPasswordForm)
	http.HandleFunc("/resetpassword", resetPassword)
	//static files public/
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	//for chrome browser to say there is not favicon
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
