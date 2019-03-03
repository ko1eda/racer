package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tinylttl/racer"
)

func main() {
	r := chi.NewRouter()
	// r.Get("/racer/chat/{chadID:[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}}")

	brokerm := make(map[string]*racer.Broker)

	r.Get("/racer/chat", serveHome)
	r.Get("/racer/cat", serveAbout)
	r.Get("/racer/chat/{chatID}", racer.ChatHandler(brokerm))

	// blocks our application
	http.ListenAndServe(":80", r)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveAbout(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "about.html")
}
