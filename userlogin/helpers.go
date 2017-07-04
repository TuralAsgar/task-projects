package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "login")
	errorHandler(w, err)
	email := session.Values["email"]
	if email == nil {
		return false
	}
	return true
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func errorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func getHashCode(data string) string {
	h := hmac.New(sha256.New, []byte("Lazy dog jumped over fire"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
