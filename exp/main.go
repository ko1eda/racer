package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tinylttl/rocket"
)

func main() {
	s := "Hello"

	fmt.Printf("%d", len(s))

	room := rocket.NewRoom()

	go room.StartBroadcast()

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		rocket.StartClientInRoom(room, w, r)
	})

	// blocks our application
	http.ListenAndServe(":80", nil)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func handle(w http.ResponseWriter, r *http.Request) {
	// check if the room is already existing

	// then pass it to the client and thats it

	// if not create a room start a broadcast and then pass it to the client
}
