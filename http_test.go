package racer_test

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/tinylttl/racer"
)

func TestHandleChat(t *testing.T) {
	cases := []struct {
		name     string
		expected *racer.Broker
	}{
		{
			name:     "Multiple requests to the same chatID endpoint should only create 1 broker",
			expected: racer.NewBroker("2"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// wrap the recorder with a Hijack method so we can use it with
			// gorillas upgrader
			w := newRecorder(*httptest.NewRecorder())
			r := reqWithSockHeaders(t, "GET", "/racer/chat/23")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("chatID", "23")
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			manager := racer.NewManager(racer.WithMap(map[string]*racer.Broker{"23": tc.expected}))
			h := racer.ChatHandler(manager)

			i := 0
			for i < 15 {
				go h(w, r)
				i++
			}

			actual, _ := manager.Exists("23")
			if actual != tc.expected {
				t.Fatalf("\tactual: %+v expected: %+v", actual, tc.expected)
			}
		})
	}

	// for _, tc := range cases {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		// wrap the recorder with a Hijack method so we can use it with
	// 		// gorillas upgrader
	// 		w := newRecorder(*httptest.NewRecorder())
	// 		r := reqWithSockHeaders(t, "GET", "/racer/chat/23")
	// 		rctx := chi.NewRouteContext()
	// 		rctx.URLParams.Add("chatID", "23")
	// 		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// 		manager := racer.NewManager()
	// 		manager.Register("23", tc.expected)
	// 		h := racer.ChatHandler(manager)

	// 		i := 0
	// 		for i < 25 {
	// 			fmt.Printf("goroutine num %d \n", i)
	// 			go h(w, r)
	// 			i++
	// 		}

	// 		actual, _ := manager.Exists("23")
	// 		if actual != tc.expected {
	// 			t.Fatalf("\tactual: %+v expected: %+v", actual, tc.expected)
	// 		}

	// 	})
	// }
}

func reqWithSockHeaders(t *testing.T, method, uri string) *http.Request {
	r := httptest.NewRequest(method, uri, nil)
	r.Header.Add("Connection", "Upgrade")
	r.Header.Add("Upgrade", "websocket")
	r.Header.Add("Sec-Websocket-Version", "13")
	r.Header.Add("Sec-WebSocket-Key", "13")

	return r
}

func newRecorder(r httptest.ResponseRecorder) *recorder {
	_, server := net.Pipe()

	return &recorder{r, server}
}

// recorder it similar to httptest.ResponseRecorder, but with Hijack capabilities
type recorder struct {
	httptest.ResponseRecorder
	server net.Conn
}

// Hijack the connection
func (r *recorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// return to the recorder the recorder, which is the recorder side of the connection
	rw := bufio.NewReadWriter(bufio.NewReader(r.server), bufio.NewWriter(r.server))
	return r.server, rw, nil
}

func NewDialer(h http.Handler) *websocket.Dialer {
	client, server := net.Pipe()
	conn := &recorder{server: server}

	// run the runServer in a goroutine, so when the Dial send the request to
	// the recorder on the connection, it will be parsed as an HTTPRequest and
	// sent to the Handler function.
	go conn.runServer(h)

	// use the websocket.NewDialer.Dial with the fake net.recorder to communicate with the recorder
	// the recorder gets the client which is the client side of the connection
	return &websocket.Dialer{NetDial: func(network, addr string) (net.Conn, error) { return client, nil }}
}

// runServer reads the request sent on the connection to the recorder
// from the websocket.NewDialer.Dial function, and pass it to the recorder.
// once this is done, the communication is done on the wsConn
func (r *recorder) runServer(h http.Handler) {
	// read from the recorder connection the request sent by the recorder.Dial,
	// and use the handler to serve this request.
	req, err := http.ReadRequest(bufio.NewReader(r.server))
	if err != nil {
		return
	}
	h.ServeHTTP(r, req)
}

// WriteHeader write HTTP header to the client and closes the connection
func (r *recorder) WriteHeader(code int) {
	resp := http.Response{StatusCode: code, Header: r.Header()}
	resp.Write(r.server)
}
