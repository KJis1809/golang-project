package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

var arrayOfPaths = []string{"/edit/", "/save/", "/view/"}
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	path := r.URL.Path

	for _, a := range arrayOfPaths {

		if strings.HasPrefix(path, a) {
			validPath := path[len(a):]

			if len(validPath) != 0 {
				return validPath, nil
			}
		}

	}
	return "", errors.New("invalid page title")
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {

	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
