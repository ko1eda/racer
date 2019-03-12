package racer_test

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/tinylttl/racer"
)

type message struct {
	Sent     time.Time `json:"sent"`
	Body     string    `json:"body"`
	SenderID int       `json:"senderID"`
}

func TestChatHandler(t *testing.T) {
	t.Run("It creates a new broker for each new chatID", func(t *testing.T) {
		manager := racer.NewManager()

		d := NewDialer(racer.ChatHandler(manager), [][]string{{"chatID", "23"}})
		d2 := NewDialer(racer.ChatHandler(manager), [][]string{{"chatID", "24"}})
		_, _, _ = d.Dial("ws://racer/chat/23", nil)
		_, _, _ = d2.Dial("ws://racer/chat/24", nil)

		want := 2
		if got := manager.Size(); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
	})

	t.Run("It removes brokers when they have no clients", func(t *testing.T) {
		manager := racer.NewManager()
		d := NewDialer(racer.ChatHandler(manager), [][]string{{"chatID", "23"}})
		d2 := NewDialer(racer.ChatHandler(manager), [][]string{{"chatID", "24"}})

		done := make(chan struct{})

		go func() {
			conn1, _, _ := d.Dial("ws://racer/chat/23", nil)
			conn1.Close()
			done <- struct{}{}
		}()

		go func() {
			conn2, _, _ := d2.Dial("ws://racer/chat/24", nil)
			conn2.Close()
			time.Sleep(200 * time.Millisecond)
			done <- struct{}{}
		}()

		var i = 2
		for i > 0 {
			<-done
			i--
		}

		var want int
		if got := manager.Size(); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
	})
}

func TestChatHandler_SocketConn(t *testing.T) {
	cases := []struct {
		name string
		want *message
	}{
		{
			name: "It translates json to a valid message",
			want: &message{Body: "Test"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manager := racer.NewManager()
			d := NewDialer(racer.ChatHandler(manager), [][]string{{"chatID", "23"}})

			conn, _, err := d.Dial("ws://racer/chat/23", nil)

			if err != nil {
				t.Fatal(err)
			}
			conn.WriteJSON(tc.want)

			var got message
			conn.ReadJSON(&got)

			if got.Body != tc.want.Body {
				t.Fatalf("got:%s , want: %s", got.Body, tc.want.Body)
			}
		})
	}
}

// borrowed and modified from https://github.com/posener/wstest
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

func NewDialer(h http.Handler, ctx ...[][]string) *websocket.Dialer {
	client, server := net.Pipe()
	conn := &recorder{server: server}

	// run the runServer in a goroutine, so when the Dial send the request to
	// the recorder on the connection, it will be parsed as an HTTPRequest and
	// sent to the Handler function.
	go conn.runServer(h, ctx...)

	// use the websocket.NewDialer.Dial with the fake net.recorder to communicate with the recorder
	// the recorder gets the client which is the client side of the connection
	return &websocket.Dialer{NetDial: func(network, addr string) (net.Conn, error) { return client, nil }}
}

// runServer reads the request sent on the connection to the recorder
// from the websocket.NewDialer.Dial function, and pass it to the recorder.
// once this is done, the communication is done on the wsConn
func (r *recorder) runServer(h http.Handler, ctx ...[][]string) {
	req, err := http.ReadRequest(bufio.NewReader(r.server))
	// fmt.Println("Called")

	if err != nil {
		return
	}

	addRouteCtx(&req, ctx[0])

	h.ServeHTTP(r, req)
}

// key value tuples [[key val]]
func addRouteCtx(req **http.Request, keyval [][]string) {
	rctx := chi.NewRouteContext()
	for _, tuple := range keyval {
		for j := 0; j < 2; j = j + 2 {
			rctx.URLParams.Add(tuple[j], tuple[j+1])
		}
	}

	*req = (*req).WithContext(context.WithValue((*req).Context(), chi.RouteCtxKey, rctx))
}

// WriteHeader write HTTP header to the client and closes the connection
func (r *recorder) WriteHeader(code int) {
	resp := http.Response{StatusCode: code, Header: r.Header()}
	resp.Write(r.server)
}
