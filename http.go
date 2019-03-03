package racer

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

// ChatHandler handles all GET requests to ../chat/{chatID}
// It takes a map that will map chatIDs to Brokers.
// The goal is that we only have one broker running for a given chat endpoint (chatID).
// The brokers job is to manage each client connection that is active at that endpoint.
// If a brokers clients all unregister, it will terminate and remove itself from the broker map.
func ChatHandler(brokerm map[string]*Broker) http.HandlerFunc {
	var mu sync.RWMutex

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chatID := chi.URLParam(r, "chatID")

		if chatID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Println("Chat Id: ", chatID)

		mu.RLock()
		b, exists := brokerm[chatID]
		mu.RUnlock()
		if !exists {
			b = NewBroker()

			mu.Lock()
			brokerm[chatID] = b
			defer mu.Unlock()

			fmt.Printf("%+v\n", brokerm)

			// Start the broker in its own go routine since it doesn't already exit.
			go func() {
				b.Start()
				mu.Lock()
				delete(brokerm, chatID)
				defer mu.Unlock()
			}()
		}

		c := NewClient()
		b.RegisterSubscriber(c)
		c.Run(w, r)
	})
}
