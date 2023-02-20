package main

import (
	"bytes"
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

func webserver(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", webHandler)
	mux.HandleFunc("/settings", settingsHandler)
	return &http.Server{Addr: addr, Handler: mux}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	index := &Page{Description: "testing"}
	render(w, "index", index)
}
func settingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	switch r.Method {
	case "GET":
		err := templates.ExecuteTemplate(w, "settings.html", currentConfig) //cannot use render() since currentConfig is not type Page
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// render(w, "settings", nil)
	case "POST":
		bodyBytes, err := io.ReadAll(r.Body) //parse request Body
		if err != nil {
			log.Println(err)
		}
		var out bytes.Buffer
		err = json.Indent(&out, bodyBytes, "", "\t") //format bytes to json(with indents) and send to out
		if err != nil {
			log.Println(err)
		}
		configFile, err := os.OpenFile("config.json", os.O_CREATE|os.O_TRUNC, 0600) //open config.json with rw access
		if err != nil {
			log.Println(err)
		}
		_, err = out.WriteTo(configFile) //replace all contents of config.json
		if err != nil {
			log.Println(err)
		}
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
