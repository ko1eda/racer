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

	go broker.StartBroadcast()
	go broker2.StartBroadcast()

	r := chi.NewRouter()

	r.Get("/", serveHome)
	r.Get("/about", serveAbout)

	r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {

		b := racer.NewBroker()

		c := &racer.Client{}

		b.BrokerClientCon(c, w, r)
	})

	r.Get("/chat123", func(w http.ResponseWriter, r *http.Request) {
		b := racer.NewBroker()

		c := &racer.Client{}

		b.BrokerClientCon(c, w, r)
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

func handle(w http.ResponseWriter, r *http.Request) {
	// check if the broker is already existing

	// then pass it to the client and thats it

	// if not create a broker start a broadcast and then pass it to the client
}
