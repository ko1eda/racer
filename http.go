package racer

import (
	"net/http"

	"github.com/go-chi/chi"
)

// ChatHandler handles all GET requests to ../chat/{chatID}
// It takes a broker manager that will map chatIDs to running brokers.
// The goal is that we only have one broker running for a given chat endpoint (chatID).
// The brokers job is to manage each client connection that is active at that endpoint.
// If a brokers clients all unregister, it will terminate and remove itself from its manager.
func ChatHandler(bm *BrokerManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chatID := chi.URLParam(r, "chatID")

		if chatID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// fmt.Println("Chat Id: ", chatID)

		bm.Lookup(chatID, func(found bool, b *Broker) {
			if !found {
				go func() {
					// fmt.Println("NOT FOUND CALLED IN HANDLER")
					b.Start() // blocking
					bm.Remove(chatID)
				}()
			}
			c := NewClient()
			b.RegisterSubscriber(c)
			c.Run(w, r)
		})
	})
}
