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

func TestHandleChat_SocketConn(t *testing.T) {
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
			d := NewDialer(racer.ChatHandler(manager))

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

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("chatID", "23")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.ServeHTTP(r, req)
}

// WriteHeader write HTTP header to the client and closes the connection
func (r *recorder) WriteHeader(code int) {
	resp := http.Response{StatusCode: code, Header: r.Header()}
	resp.Write(r.server)
}
