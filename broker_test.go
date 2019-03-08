package racer_test

import (
	"sync"
	"testing"

	"github.com/tinylttl/racer"
)

// TestLookup tests the basic functionality of the BrokerManagers Lookup method
func TestLookup(t *testing.T) {
	cases := []struct {
		name       string
		key        string
		brokerm    map[string]*racer.Broker
		wantBroker *racer.Broker
		wantFound  bool
	}{
		{name: "It reuses an existing broker if it has an entry in the map", wantFound: true, key: "23", brokerm: map[string]*racer.Broker{"23": racer.NewBroker("23")}},
		{name: "It creates a new broker if it cannot find one", wantFound: false, key: "24"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testmap := make(map[string]*racer.Broker)

			if len(tc.brokerm) > 0 {
				t.Logf("Map provided, Using tc.brokerm as testmap")
				testmap = tc.brokerm
			}

			bm := racer.NewManager(racer.WithMap(testmap))
			tc.wantBroker = testmap[tc.key]

			bm.Lookup(tc.key, func(found bool, b *racer.Broker) {
				// determines if the broker was found or not in the map
				if found != tc.wantFound {
					t.Fatalf("got: %+v, want: %+v", found, tc.wantFound)
				}

				// if our test map had no existing broker, then that means we expect one to be created
				if tc.wantBroker == nil {
					tc.wantBroker = b
				}

				if b != tc.wantBroker {
					t.Fatalf("got: %+v, want: %+v", b, tc.wantBroker)
				}
			})
		})
	}
}

// TestLookupConcurrent ensures that the number of lookups stays consistent across multiple go routines,
// It tests that only one new broker is ever created if one does not exist for a given key.
func TestLookupConcurrent(t *testing.T) {
	cases := []struct {
		name  string
		keys  []string
		count chan struct{}
		want  int
	}{
		{name: "It sets found to false only once if no broker is found for a key", keys: []string{"23"}, count: make(chan struct{}, 15), want: 1},
		{name: "It creates a new broker for each new key it encounters", keys: []string{"23", "24"}, count: make(chan struct{}, 15), want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bm := racer.NewManager()

			var wg sync.WaitGroup
			wg.Add(1)
			i := 0
			// this will create 10 * n gorotuines where n is the number of keys in tc.keys
			for i < 10 {
				for _, key := range tc.keys {
					wg.Add(1)
					go bm.Lookup(key, func(found bool, b *racer.Broker) {
						if found == false {
							tc.count <- struct{}{}
						}
						wg.Done()
					})
				}
				i++
			}
			wg.Done() // this ensures we wait until the loop finishes and remove the wg.Add(1) before the loop started
			wg.Wait()
			close(tc.count)

			got := 0
			for _ = range tc.count {
				got++
			}

			if got != tc.want {
				t.Fatalf("got: %d, want: %d", got, tc.want)
			}
		})
	}
}

// TestRemoveConcurrent tests that multiple concurrent calls to
// a managers remove function can be made safely.
// These tests cover all code in the Remove function
func TestRemoveConcurrent(t *testing.T) {
	cases := []struct {
		testm map[string]*racer.Broker
		name  string
		keys  []string // to remove
		want  int
	}{
		{
			name: "It removes multiple brokers at once",
			keys: []string{"10291", "191", "1589Adx1"},
			testm: map[string]*racer.Broker{
				"10291":    racer.NewBroker("10291"),
				"xx90":     racer.NewBroker("xx90"),
				"191":      racer.NewBroker("191"),
				"12":       racer.NewBroker("12"),
				"1589Adx1": racer.NewBroker("1589Adx1"),
			},
			want: 2,
		},
		{
			name: "It handles removal of non-existant keys",
			keys: []string{"10291", "191", "1589Adx1"},
			testm: map[string]*racer.Broker{
				"11111": racer.NewBroker("11111"),
			},
			want: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := len(tc.testm)
			manager := racer.NewManager(racer.WithMap(tc.testm))

			i := 0
			done := make(chan struct{}, len(tc.keys))
			for i < len(tc.keys) {
				// one go routine per removal simulates one broker per chatID simulatenously
				// removing itself from the manager
				// EXAMPLE: a scenario where mulltiple different chats have all their clients log out at the same time
				// this should not cause an issue with the manager.
				go func(i int) {
					if removed := manager.Remove(tc.keys[i]); !removed {
						if got != tc.want {
							t.Fatalf("got: %d, want: %d", got, tc.want)
						}
					}
					done <- struct{}{}
				}(i)

				i++
			}

			d := len(tc.keys)
			for d > 0 {
				<-done
				d--
			}

			// the final size of the map should be
			// the intial size of map - num keys removed
			got = len(tc.testm)
			if got != tc.want {
				t.Fatalf("got: %d, want: %d", got, tc.want)
			}
		})
	}
}

type subscriber interface {
	Register(broadcast chan<- []byte, unregister chan chan<- []byte) (send chan<- []byte)
}

type client struct {
	broadcast  chan<- []byte
	unregister chan chan<- []byte
	send       chan []byte
}

func (c *client) Register(broadcast chan<- []byte, unregister chan chan<- []byte) chan<- []byte {
	c.broadcast = broadcast
	c.unregister = unregister
	c.send = make(chan []byte)
	return c.send
}

func TestStart(t *testing.T) {
	// we have a broker it runs
	cases := []struct {
		name   string
		want   []byte
		broker *racer.Broker
	}{
		{
			name:   "It registers subscribers",
			broker: racer.NewBroker("x"),
			want:   []byte("Test Message"),
		},
		{
			name:   "It unregisters subscribers",
			broker: racer.NewBroker("x"),
			want:   make([]byte, 1),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sub := &client{}

			go func() {
				tc.broker.Start()
			}()

			tc.broker.RegisterSubscriber(sub)

			sub.broadcast <- tc.want

			if got := <-sub.send; string(got) != string(tc.want) {
				t.Fatalf("got: %s, want: %s", string(got), tc.want)
			}

			// if we get something other than a close signal on the chan we have a problem
			sub.unregister <- sub.send
			if got, ok := <-sub.send; ok {
				t.Fatalf("got: %s, want: %s", string(got), string(tc.want))
			}
		})
	}
	// we want to register a new client

	// we want to unregister a client

	// we want to

}
