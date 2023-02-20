package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Page struct {
	Title       string
	Description string
	Body        []byte
}

func webserver(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", webHandler)
	mux.HandleFunc("/settings", settingsHandler)
	return &http.Server{Addr: addr, Handler: mux}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "hello world")
	index := &Page{Description: "testing"}
	render(w, "index", index)
}
func settingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	switch r.Method {
	case "GET":
		render(w, "settings", nil)
	case "POST":
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(bytes))
		configFile, err := os.OpenFile("config.json", os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Println(err)
		}
		configFile.WriteString(string(bytes))
		err = configFile.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

var templates = template.Must(template.ParseFiles("./web/index.html", "./web/settings.html"))

func render(w http.ResponseWriter, template string, p *Page) {
	err := templates.ExecuteTemplate(w, template+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
