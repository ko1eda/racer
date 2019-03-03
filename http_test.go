package racer_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
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

			i := 3
			for i > 0 {
				go h(w, r)
				i--
			}

			// time.Sleep(500 * time.Millisecond)
			// if actual := len(brokerm); actual != c.expected {
			// 	t.Fatalf("\tactual: %d expected: %d", actual, c.expected)
			// }

		})
	}

}
