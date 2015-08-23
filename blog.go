package main

import (
	"io/ioutil"
	"net/http"
)

// Post stores the data of a single blog post
type Post struct {
	Title string
	Body  []byte
}

func (p *Post) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPost(title string) (*Post, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Post{Title: title, Body: body}, nil

}

func main() {
	http.HandleFunc("/admin/login/auth", authHandler)
	http.HandleFunc("/admin/login", loginHandler)
	http.HandleFunc("/admin/logout", logoutHandler)
	http.HandleFunc("/admin/", homeHandler)

	http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
}
