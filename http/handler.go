package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tinylttl/racer"
	"github.com/tinylttl/racer/broker"
	"github.com/tinylttl/racer/gorilla"
)

// Handler handles all incoming HTTP requests for the application
type Handler struct {
	Router chi.Router
	Repo   racer.MessageRepo
}

// NewHandler returns a Handler configured with a Router.
func NewHandler(repo racer.MessageRepo) *Handler {
	h := &Handler{Repo: repo}

	h.Router = NewRouter(h)

	return h
}

// NewRouter returns a new router preloaded with all the routes necessary to serve
// the application.
func NewRouter(handler *Handler) chi.Router {
	routeBase := "/v" + apiVersion

	r := chi.NewRouter()
	b := broker.NewBroker()

	r.Get(routeBase+"/chat/{chatID}", handler.handleGetTopic(b))

	return r
}

// ServeHTTP wraps the routers own ServeHTTP method
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

// handleGetTopic handles all GET requests to /chat/:chatID
// It takes a broker that maps IDS to running topics.
// The goal is that we only have one topic running for a given chat endpoint (chatID).
// The topics job is to manage each client connection that is active at that endpoint.
// If a topics clients all unregister, it will terminate and remove itself from the broker.
func (h *Handler) handleGetTopic(b *broker.Broker) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chatID := chi.URLParam(r, "chatID")

		if chatID == "" {
			// Log here since this is where we handle the error
			http.Error(w, "", http.StatusNotFound)
			return
		}

		b.Lookup(chatID, func(found bool, t *broker.Topic) {
			if !found {
				go func() {
					t.Start() // TODO: PASS CONTEXT TO CANCEL
					b.Remove(chatID)
				}()
			}

			conn, err := gorilla.NewConnection(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) 
				return
			}

			backupper := racer.NewBackupper(chatID, h.Repo)
			
			c := racer.NewClient(t, conn, backupper)
			c.Run()
		})
	})
}

// r.Get("/racer/chat/{chadID:[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}}")
