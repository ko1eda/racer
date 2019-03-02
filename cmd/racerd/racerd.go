package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tinylttl/racer"
)

func main() {
	r := chi.NewRouter()
	// r.Get("/racer/chat/{chadID:[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}}")

	brokerm := make(map[string]*racer.Broker)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the chatid from the uri
		chatID := chi.URLParam(r, "chatID")
		if chatID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Println("Chat Id: ", chatID)
		// check if the broker is in our map of known brokers
		// 1 broker per 1 chat uri
		b, ok := brokerm[chatID]
		c := racer.NewChat()
		if !ok {
			b = racer.NewBroker()

			brokerm[chatID] = b

			fmt.Printf("%+v\n", brokerm)

			// Start the broker in its own go routine
			// TODO add context on all the clients that is tied to the brokers parent context
			// if the broker fails, all the clients should gracefully stop and save any messages
			// Make sure theres no race condition writing to this shared map
			go func() {
				defer delete(brokerm, chatID)
				b.Start()
			}()
		}

		// register the new client with the broker
		// and then run the clients dameon process which will read and write to a socket
		// with messages broadcast through its broker
		b.RegisterSubscriber(c)
		c.Run(w, r)

		fmt.Println("HANDLER FUNC ENDED")
	})

	r.Get("/racer/chat", serveHome)
	r.Get("/racer/scat", serveAbout)
	r.Get("/racer/chat/{chatID}", h)

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
