package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

// Get /
func index(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	dashboard(w, r)
}

//function used in index
func dashboard(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	session, err := store.Get(r, "login")
	email := session.Values["email"]
	rows, err := db.Query(`SELECT Id, Email, Name, Address, Phone FROM user WHERE email=?`, email)
	check(err)

	User := struct {
		Id, Email, Name, Address, Phone string
	}{}

	if rows.Next() {
		err = rows.Scan(&User.Id, &User.Email, &User.Name, &User.Address, &User.Phone)
		check(err)
	}
	tpl.ExecuteTemplate(w, "dashboard.html", User)
}

func signupform(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	repeatPassword := r.FormValue("repeatPassword")

	if password != repeatPassword && password != "" {
		w.Write([]byte("The pasword fields don't match or empty!"))
		return
	}

	rows, err := db.Query(`SELECT * FROM user WHERE email=?`, email)
	check(err)

	if !rows.Next() {
		_, err = db.Exec(`INSERT INTO user (email, password) VALUES (?,?)`, email, getHashCode(password))
		check(err)
		session, err := store.Get(r, "login")
		errorHandler(w, err)

		session.Values["email"] = email
		session.Save(r, w)

		w.Write([]byte("You have logged in successfully"))
		return
	}
	http.Redirect(w, r, "Sorry this email exists", http.StatusConflict)

}

func completeform(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "complete.html", nil)
}

func complete(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "login")
	errorHandler(w, err)
	email := session.Values["email"]
	if email == "" || email == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	name := r.FormValue("name")
	address := r.FormValue("address")
	phone := r.FormValue("phone")
	if address != "" && phone != "" && name != "" {
		_, err = db.Exec(`UPDATE user SET name=?,address=?,phone=? WHERE email=?`,
			name, address, phone, email)
		check(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Fprint(w, "Fill out all the fields")

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email != "" && password != "" {
		password = getHashCode(password)
		rows, err := db.Query(`SELECT * FROM user WHERE email=? and password=?`, email, password)
		check(err)

		user := struct {
			id, email, password, name, address, phone string
		}{}

		if !rows.Next() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Something went wrong!")
			return
		}
		err = rows.Scan(&user.id, &user.email, &user.password, &user.name, &user.address, &user.phone)
		check(err)

		session, err := store.Get(r, "login")
		errorHandler(w, err)

		session.Values["username"] = user.name
		session.Values["email"] = user.email
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
		w.Write([]byte("You have logged in successfully"))
		return
	}
	w.Write([]byte("Fill out all fields"))

}

func logingoogle(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("google")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
			http.StatusBadRequest)
		return
	}
	loginUrl, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err),
			http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, loginUrl, http.StatusTemporaryRedirect)

}

func callback(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("google")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",
			provider, err), http.StatusBadRequest)
		return
	}
	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err),
			http.StatusInternalServerError)
		return
	}
	user, err := provider.GetUser(creds)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s",
			provider, err), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`SELECT * FROM user WHERE email=?`, user.Email())
	check(err)
	session, err := store.Get(r, "login")
	errorHandler(w, err)

	if !rows.Next() {
		_, err = db.Exec(`INSERT INTO user (email) VALUES (?)`, user.Email())
		check(err)
		errorHandler(w, err)

		session.Values["email"] = user.Email()
		session.Save(r, w)

		http.Redirect(w, r, "/completeform", http.StatusFound)
		return
	}

	session.Values["email"] = user.Email()
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("login")
	c = &http.Cookie{
		Name:   "login",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func editform(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	tpl.ExecuteTemplate(w, "edit.html", nil)

}

func edit(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	//GET
	if r.Method == http.MethodGet {
		session, err := store.Get(r, "login")
		errorHandler(w, err)
		email := session.Values["email"]
		rows, err := db.Query(`SELECT id, email, name, address, phone  FROM user WHERE email=?`, email)
		check(err)

		User := struct {
			Id, Email, Name, Address, Phone string
		}{}

		if rows.Next() {
			err = rows.Scan(&User.Id, &User.Email, &User.Name, &User.Address, &User.Phone)
			check(err)
		}
		json.NewEncoder(w).Encode(User)
		return
	}
	//POST
	session, err := store.Get(r, "login")
	errorHandler(w, err)
	email := session.Values["email"]
	name := r.FormValue("name")
	password := r.FormValue("password")
	repeatPassword := r.FormValue("repeatPassword")
	address := r.FormValue("address")
	phone := r.FormValue("phone")
	if password == repeatPassword && (password != "" && repeatPassword != "" && address != "" && phone != "" && name != "") {
		_, err = db.Exec(`UPDATE user SET name=?,password=?,address=?,phone=?,email=? WHERE email=?`,
			name, getHashCode(password), address, phone, email, email)
		check(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Fprint(w, "The password you typed doesn't fit")

}

func forgot(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	tpl.ExecuteTemplate(w, "forgot.html", nil)

}

func reset(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email != "" {
		rows, err := db.Query(`SELECT * FROM user WHERE email=?`, email)
		check(err)

		if !rows.Next() {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "This email doesn't exist")
			return
		}

		rows, err = db.Query(`SELECT * FROM reset WHERE email=?`, email)
		check(err)

		if rows.Next() {
			_, err = db.Exec(`DELETE FROM reset WHERE email=?`, email)
			check(err)
		}

		random, _ := Random(randomBytes)
		_, err = db.Exec(`INSERT INTO reset (email,random) VALUES (?,?)`, email, random)
		check(err)
		sender := NewSender("user", "password")

		//The receiver needs to be in slice as the receive supports multiple receiver
		Receiver := []string{email}

		Subject := "Reset password"

		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

		bodyMessage := fmt.Sprintf(
			"%s<html><body><a href='http://localhost:8080/resetpasswordform?email=%s&val=%s'>Click here to reset your password</a></body></html>",
			mime, email, random)

		sender.SendMail(Receiver, Subject, bodyMessage)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func resetPassword(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	secret := r.FormValue("val")
	password := r.FormValue("password")
	repeatPassword := r.FormValue("repeatPassword")

	rows, err := db.Query(`SELECT * FROM user WHERE email=?`, email)
	check(err)

	if !rows.Next() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "This email doesn't exist")
		return
	}

	rows, err = db.Query(`SELECT random FROM reset WHERE email=?`, email)
	check(err)

	var random string

	if rows.Next() {
		err = rows.Scan(&random)
		check(err)
	}

	if password == repeatPassword && password != "" {
		if secret == random {
			_, err = db.Exec(`UPDATE user SET password=? WHERE email=?`, getHashCode(password), email)
			check(err)
			_, err = db.Exec(`DELETE FROM reset WHERE email=?`, email)
			check(err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		fmt.Fprint(w, "The secret code is not correct")
		return
	}
	fmt.Fprint(w, "Passwords can't be empty and must match")

}

func resetPasswordForm(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	secret := r.FormValue("val")

	rows, err := db.Query(`SELECT * FROM user WHERE email=?`, email)
	check(err)

	if !rows.Next() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "This email doesn't exist")
		return
	}

	rows, err = db.Query(`SELECT random FROM reset WHERE email=?`, email)
	check(err)

	var random string

	if rows.Next() {
		err = rows.Scan(&random)
		check(err)
	}

	if secret == random {
		tpl.ExecuteTemplate(w, "reset.html", nil)
		return
	}
	fmt.Fprint(w, "The secret code is not correct")
	return

}
