package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/tinylttl/rocket"
)

func main() {
	room := rocket.NewRoom()

	room2 := rocket.NewRoom()

	go room.StartBroadcast()
	go room2.StartBroadcast()

	r := chi.NewRouter()

	r.Get("/", serveHome)
	r.Get("/about", serveAbout)

	r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(10)
		rocket.StartClientInRoom(id, room, w, r)
	})

	r.Get("/chat123", func(w http.ResponseWriter, r *http.Request) {
		rocket.StartClientInRoom(2000, room2, w, r)
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
	// check if the room is already existing

	// then pass it to the client and thats it

	// if not create a room start a broadcast and then pass it to the client
}
