package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

// Struct to hold individual messages with timestamps
type Message struct {
	Text string
	Time string
}

// Struct to hold all messages
type PageData struct {
	Messages []Message
}

// Mutex to handle concurrent writes safely
var messageMutex sync.Mutex
var messages []Message

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
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

	err = tmplParsed.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	messageMutex.Lock()
	defer messageMutex.Unlock()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Create new message with timestamp
		newMessage := Message{
			Text: r.FormValue("message"),
			Time: time.Now().Format("2006-01-02 15:04:05"),
		}

		// Append to the message history
		messages = append(messages, newMessage)
	}

	// Pass the message history to the template
	data := PageData{Messages: messages}
	renderTemplate(w, "content", data)
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
