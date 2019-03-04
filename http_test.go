package racer_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/tinylttl/racer"
)

func TestHandleChat(t *testing.T) {
	cases := []struct {
		name     string
		expected int
	}{
		{
			name:     "Multiple clients connecting to the same chat broker should produce one map entry",
			expected: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := newRecorder(*httptest.NewRecorder())
			r := httptest.NewRequest("GET", "/racer/chat/23", nil)
			r.Header.Add("Connection", "Upgrade")
			r.Header.Add("Upgrade", "websocket")
			r.Header.Add("Sec-Websocket-Version", "13")
			r.Header.Add("Sec-WebSocket-Key", "13")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("chatID", "23")
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			brokerm := racer.NewManager()

			h := racer.ChatHandler(brokerm)

			i := 0
			for i < 25 {
				fmt.Printf("goroutine num %d \n", i)
				go h(w, r)
				// time.Sleep(2000 * time.Millisecond)
				i++
			}

			time.Sleep(500 * time.Millisecond)
			_, actual := brokerm.BrokerExists("23")
			if actual != true {
				t.Fatalf("\tactual: %v expected: %d", actual, tc.expected)
			}

		})
	}
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

func newRecorder(r httptest.ResponseRecorder) *recorder {
	_, server := net.Pipe()

	return &recorder{r, server}
}

// recorder it similar to httptest.ResponseRecorder, but with Hijack capabilities
type recorder struct {
	httptest.ResponseRecorder
	server net.Conn
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

// Hijack the connection
func (r *recorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// return to the recorder the recorder, which is the recorder side of the connection
	rw := bufio.NewReadWriter(bufio.NewReader(r.server), bufio.NewWriter(r.server))
	return r.server, rw, nil
}

// WriteHeader write HTTP header to the client and closes the connection
func (r *recorder) WriteHeader(code int) {
	resp := http.Response{StatusCode: code, Header: r.Header()}
	resp.Write(r.server)
}

// for _, tc := range cases {
// 	t.Run(tc.name, func(t *testing.T) {
// 		brokerm := racer.NewManager()
// 		r := chi.NewRouter()
// 		r.Get("/racer/chat/{chatID}", racer.ChatHandler(brokerm))

// 		s := httptest.NewServer(r)
// 		defer s.Close()

// 		u := url.URL{Scheme: "ws", Host: "ws" + strings.TrimPrefix(s.URL, "http"), Path: "/racer/chat/23"}
// 		i := 1
// 		for i > 0 {
// 			go func() {
// 				client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				time.Sleep(500 * time.Millisecond)

// 				client.WriteMessage(websocket.TextMessage, []byte("HELLO"))

// 				time.Sleep(500 * time.Millisecond)

// 				// client.Close()
// 			}()
// 			i--
// 		}

// 		time.Sleep(2 * time.Second)

// 		log.Fatal(brokerm.BrokerExists("23"))
// 		// ln, err := net.Listen("tcp", "127.0.0.1:0")
// 		// var server net.Conn
// 		// go func() {
// 		// 	for {
// 		// 		defer server.Close()
// 		// 		server, err = ln.Accept()
// 		// 		if err != nil {
// 		// 			return
// 		// 		}
// 		// 	}
// 		// 	http.Serve(ln, racer.ChatHandler(brokerm))
// 		// }()
// 	})
// }

// func (t *testing.T) startServer()
