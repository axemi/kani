package main

import (
	"net/http"
	"text/template"
)

type Page struct {
	Title       string
	Description string
	Body        []byte
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "hello world")
	index := &Page{Description: "testing"}
	render(w, "index", index)
}

var templates = template.Must(template.ParseFiles("./web/index.html"))

func render(w http.ResponseWriter, template string, p *Page) {
	err := templates.ExecuteTemplate(w, template+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func webserver(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", webHandler)
	return &http.Server{Addr: addr, Handler: mux}
}
