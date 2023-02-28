package main

import (
	"encoding/json"
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

type WebServer struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewWebServer(addr string) *WebServer {
	mux := http.NewServeMux()
	server := &http.Server{Addr: addr, Handler: mux}
	return &WebServer{server, mux}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	index := &Page{Description: "testing"}
	render(w, "index", index)
}

func settingsHandler(c *config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		switch r.Method {
		case "GET":
			err := templates.ExecuteTemplate(w, "settings.html", c) //cannot use render() since currentConfig is not type Page
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			// render(w, "settings", nil)
		case "POST":
			bodyBytes, err := io.ReadAll(r.Body) //parse request Body
			if err != nil {
				log.Println(err)
			}
			newConfig := &config{}
			err = json.Unmarshal(bodyBytes, newConfig) //unmarshal into newConfig so we only take data we want
			if err != nil {
				log.Println(err)
			}
			ncBytes, err := json.MarshalIndent(newConfig, "", "\t")
			if err != nil {
				log.Println(err)
			}
			configFile, err := os.OpenFile("config.json", os.O_CREATE|os.O_TRUNC, 0600) //open config.json with rw access
			if err != nil {
				log.Println(err)
			}
			_, err = configFile.Write(ncBytes) //write new config to file
			if err != nil {
				log.Println(err)
			}
		}
	})
}

var templates = template.Must(template.ParseFiles("./web/index.html", "./web/settings.html"))

func render(w http.ResponseWriter, template string, p *Page) {
	err := templates.ExecuteTemplate(w, template+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
