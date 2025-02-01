package main

import (
	"html/template"
	"log"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	tmplFiles := []string{
		"./cmd/web/templates/base.layout.html",
		"./cmd/web/templates/header.partial.html",
		"./cmd/web/templates/footer.partial.html",
		"./cmd/web/templates/test.page.html",
	}

	tmplParsed, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	err = tmplParsed.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "content")
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./cmd/web/static"))))

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
