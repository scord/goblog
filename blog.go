package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

// Post stores the data of a single blog post
type Post struct {
	Title string
	Body  []byte
}

func (p *Post) saveToFile() error {
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

func getPosts(db *sql.DB) ([]string, error) {
	rows, _ := db.Query("SELECT title FROM posts")
	var ps []string
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./blog.db")
	posts, _ := getPosts(db)
	renderPostsTemplate(w, posts)
}

func (p *Post) save(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO posts(id, title, body) values (?, ?, ?)", nil, p.Title, p.Body)
	return err
}

func renderPostsTemplate(w http.ResponseWriter, titles []string) {
	t, _ := template.ParseFiles("posts.html")
	t.Execute(w, titles)
}

func newpostHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "newpost")
}

func addpostHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./blog.db")
	newPost := &Post{Title: r.PostFormValue("title"), Body: []byte(r.PostFormValue("body"))}
	newPost.save(db)
	http.Redirect(w, r, "/admin", 302)
}

func main() {

	http.HandleFunc("/admin/login/auth", authHandler)
	http.HandleFunc("/admin/login", loginHandler)
	http.HandleFunc("/admin/logout", logoutHandler)
	http.HandleFunc("/admin/", homeHandler)
	http.HandleFunc("/admin/newpost", newpostHandler)
	http.HandleFunc("/admin/addpost", addpostHandler)
	http.HandleFunc("/admin/posts", postsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
}
