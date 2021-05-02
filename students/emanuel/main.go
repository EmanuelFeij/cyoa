package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	MyJSON        map[string]Page
	templatesHtml *template.Template
)

type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
type Page struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []Options `json:"options"`
}

func LoadJSON(filename string) map[string]Page {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicln("cannot open file")
	}
	decoder := json.NewDecoder(file)
	v := make(map[string]Page)
	err = decoder.Decode(&v)
	if err != nil {
		log.Panicln("cannot decode JSON")
	}

	return v

}

type MyHandler struct{}

func (MyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	if path == "/" {
		path = "intro"
	} else {
		path = path[1:]
	}
	page, ok := MyJSON[path]

	if !ok {
		http.Redirect(rw, r, "/", http.StatusNotFound)
		return
	}

	err := templatesHtml.ExecuteTemplate(rw, "index.html", &page)
	if err != nil {
		log.Panicln(err)
	}

}

func MyServeMux() *http.ServeMux {
	return http.NewServeMux()
}

func main() {

	MyJSON = LoadJSON("gopher.json")
	templatesHtml = template.Must(template.ParseFiles("templates/index.html"))

	mHandler := MyHandler{}

	myServeMux := MyServeMux()
	myServeMux.Handle("/", mHandler)
	log.Println("Serving on http://localhost:8080")
	http.ListenAndServe(":8080", myServeMux)

}
