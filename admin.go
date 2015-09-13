package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/securecookie"

	"golang.org/x/crypto/bcrypt"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func authHandler(w http.ResponseWriter, r *http.Request) {
	password := []byte(r.PostFormValue("password"))
	username := r.PostFormValue("username")

	file, err := ioutil.ReadFile("users.json")
	var users map[string]string
	json.Unmarshal(file, &users)
	fmt.Print(users[username])
	hash := []byte(users[username])
	redirect := "/admin/login"
	err = bcrypt.CompareHashAndPassword(hash, password)
	if err == nil {
		setCookie(username, w)
		redirect = "/admin"
	}
	http.Redirect(w, r, redirect, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearCookie(w)
	http.Redirect(w, r, "/admin", 302)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login")
}

func checkCookie(request *http.Request) (username string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["username"]
		}
	}
	return username
}

func setCookie(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/admin/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/admin/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func renderTemplate(w http.ResponseWriter, name string) {
	t, _ := template.ParseFiles(name + ".html")
	t.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	username := checkCookie(r)
	if username != "" {
		renderTemplate(w, "home")
	} else {
		http.Redirect(w, r, "/admin/login", 302)
	}
}

func createUser(username string, password string) error {

	jsondata, _ := ioutil.ReadFile("users.json")
	var users map[string]string
	json.Unmarshal(jsondata, &users)

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	users[username] = string(hash)

	jsondata, err = json.Marshal(users)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("users.json", jsondata, 0600)
	if err != nil {
		return err
	}

	return nil
}
