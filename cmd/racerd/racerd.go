package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tinylttl/racer"
)

func main() {
	broker := racer.NewBroker()

	broker2 := racer.NewBroker()

	go broker.Start()
	go broker2.Start()

	r := chi.NewRouter()

	r.Get("/", serveHome)
	r.Get("/about", serveAbout)

	r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
		c := racer.NewChat()

		broker.RegisterSubscriber(c)

		c.Run(w, r)
	})

	r.Get("/chat123", func(w http.ResponseWriter, r *http.Request) {
		c := racer.NewChat()

		broker2.RegisterSubscriber(c)

		c.Run(w, r)
	})

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
