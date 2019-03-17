package broker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/tinylttl/racer/broker"
)

// TestLookup tests the basic functionality of the BrokerManagers Lookup method
func TestLookup(t *testing.T) {
	cases := []struct {
		name       string
		key        string
		topicm     map[string]*broker.Topic
		wantBroker *broker.Topic
		wantFound  bool
	}{
		{name: "It reuses an existing topic if it has an entry in the map", wantFound: true, key: "23", topicm: map[string]*broker.Topic{"23": broker.NewTopic("23")}},
		{name: "It creates a new topic if it cannot find one", wantFound: false, key: "24"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testmap := make(map[string]*broker.Topic)

			if len(tc.topicm) > 0 {
				t.Logf("Map provided, Using tc.topicm as testmap")
				testmap = tc.topicm
			}

			bm := broker.NewBroker(broker.WithMap(testmap))
			tc.wantBroker = testmap[tc.key]

			bm.Lookup(tc.key, func(found bool, b *broker.Topic) {
				// determines if the topic was found or not in the map
				if found != tc.wantFound {
					t.Fatalf("got: %+v, want: %+v", found, tc.wantFound)
				}

				// if our test map had no existing topic, then that means we expect one to be created
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
// It tests that only one new topic is ever created if one does not exist for a given key.
func TestLookupConcurrent(t *testing.T) {
	cases := []struct {
		name  string
		keys  []string
		count chan struct{}
		want  int
	}{
		{name: "It sets found to false only once if no topic is found for a key", keys: []string{"23"}, count: make(chan struct{}, 15), want: 1},
		{name: "It creates a new topic for each new key it encounters", keys: []string{"23", "24"}, count: make(chan struct{}, 15), want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bm := broker.NewBroker()

			var wg sync.WaitGroup
			wg.Add(1)
			i := 0
			// this will create 10 * n gorotuines where n is the number of keys in tc.keys
			for i < 10 {
				for _, key := range tc.keys {
					wg.Add(1)
					go bm.Lookup(key, func(found bool, b *broker.Topic) {
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
		testm map[string]*broker.Topic
		name  string
		keys  []string // to remove
		want  int
	}{
		{
			name: "It removes multiple brokers at once",
			keys: []string{"10291", "191", "1589Adx1"},
			testm: map[string]*broker.Topic{
				"10291":    broker.NewTopic("10291"),
				"xx90":     broker.NewTopic("xx90"),
				"191":      broker.NewTopic("191"),
				"12":       broker.NewTopic("12"),
				"1589Adx1": broker.NewTopic("1589Adx1"),
			},
			want: 2,
		},
		{
			name: "It handles removal of non-existant keys",
			keys: []string{"10291", "191", "1589Adx1"},
			testm: map[string]*broker.Topic{
				"11111": broker.NewTopic("11111"),
			},
			want: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := len(tc.testm)
			manager := broker.NewBroker(broker.WithMap(tc.testm))

			i := 0
			done := make(chan struct{}, len(tc.keys))
			for i < len(tc.keys) {
				// one go routine per removal simulates one topic per chatID simulatenously
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

type client struct {
	broadcast  chan<- *broker.Message
	unregister chan chan<- *broker.Message
	send       chan *broker.Message
}

func (c *client) Register(broadcast chan<- *broker.Message, unregister chan chan<- *broker.Message) chan<- *broker.Message {
	c.broadcast = broadcast
	c.unregister = unregister
	return c.send
}

func TestStart(t *testing.T) {
	// we have a topic it runs
	cases := []struct {
		name  string
		want  []byte
		topic *broker.Topic
	}{
		{
			name:  "It registers subscribers",
			topic: broker.NewTopic("x"),
			want:  []byte("Test Message"),
		},
		{
			name:  "It unregisters subscribers",
			topic: broker.NewTopic("x"),
			want:  make([]byte, 1),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sub := &client{send: make(chan *broker.Message, 1)}
			sub2 := &client{send: make(chan *broker.Message, 1)}
			closed := make(chan struct{})

			go func() {
				tc.topic.Start()
				close(closed)
			}()

			tc.topic.RegisterSubscriber(sub)
			tc.topic.RegisterSubscriber(sub2)

			// A buffered channel is necessary based on the way we have the select statement in
			// topic set up, without a buffered channel any blocking for any reason on a send channel will
			// cause a close on the clients channel. A buffered channel gives us a threshold to say
			// we want to accept x amount of requests and if we block, then we know there's a problem
			// NOTE: to check that both of these recieved the same message without using a buffer
			// we would need to put one in its own goroutine, because the first receive would block the second
			// channel from recieveing (they both need to be able to recieve from the topic at the same time when a message is broadcast)
			go func() { sub.broadcast <- &broker.Message{Payload: tc.want} }()

			if got := <-sub.send; string(got.Payload.([]byte)) != string(tc.want) {
				t.Fatalf("got: %s, want: %s", string(got.Payload.([]byte)), tc.want)
			}

			// check that sub2 got the same
			if got := <-sub2.send; string(got.Payload.([]byte)) != string(tc.want) {
				t.Fatalf("got: %s, want: %s", string(got.Payload.([]byte)), tc.want)
			}

			// if we get something other than a close signal on the chan we have a problem
			sub.unregister <- sub.send
			sub2.unregister <- sub2.send

			if got, ok := <-sub.send; ok {
				t.Fatalf("got: %s, want: %s", string(got.Payload.([]byte)), string(tc.want))
			}

			// give the topic some time to unregister both channels
			time.Sleep(10 * time.Millisecond)

			// The topic should break out its start() process and close then chan when all clients unregister
			// if not we have a problem
			select {
			case <-closed:
			default:
				t.Fatalf("Topic never stopped running")
			}
		})
	}
}
