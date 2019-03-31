package racer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tinylttl/racer/broker"
)

// Client represents that chat client. Everytime a new connection is made
// to the server, a new client is created.
type Client struct {
	Broadcaster Broadcaster
	Conn        Connector
	Receive     chan *broker.Message // receive messages from the broadcaster
	Backupper   *Backupper
	ID          string
}

// Connector is the source of data to and from the client and server.
// The default connection type for racer is socket.
type Connector interface {
	Read() <-chan *Message
	Write() chan<- *Message
}

// Broadcaster can broadcast messages to other listening client goroutines
type Broadcaster interface {
	Register() chan chan<- *broker.Message // switch this back to the old register method approach with subscriber Register(*Client)
	Unregister() chan chan<- *broker.Message
	Broadcast() chan<- *broker.Message
}

// NewClient returns a new Chat client instance that is registered with a broadcaster
func NewClient(b Broadcaster, conn Connector) *Client {
	c := &Client{ID: fmt.Sprintf("%d", rand.Intn(100000)), Receive: make(chan *broker.Message, 1), Broadcaster: b, Conn: conn}

	c.Broadcaster.Register() <- c.Receive

	return c
}

// Run starts three goroutines. The first backs up messages to an in-memory data store at a set interval.
// The second reads incoming messages from the Clients connection and broadcasts them to all other clients sharing the same broadcaster.
// The third reads messages recieved from said broadcaster finally writing them back through to the connection.
func (c *Client) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	// TODO: Possibly switch these to prioritize select statements
	// to handle case where backupper returns an error
	// https://stackoverflow.com/questions/46200343/force-priority-of-go-select-statement/46202533#46202533
	go func() {
		c.Backupper.Run(ctx)
	}()

	go func() {
		for msg := range c.Conn.Read() {
			c.Broadcaster.Broadcast() <- &broker.Message{Payload: msg}

			c.Backupper.Hold(msg)
		}
		c.Broadcaster.Unregister() <- c.Receive // shutdown the client because the connection was closed
	}()

	go func() {
		for bmsg := range c.Receive {
			c.Conn.Write() <- bmsg.Payload.(*Message)
		}
		cancel()
	}()
}

// Message is data that is sent as json through the connection.
type Message struct {
	Timestamp int64  `json:"timestamp"`
	Sent      string `json:"sent"`
	Body      string `json:"body"`
	SenderID  int    `json:"senderID"`
}

// MessageRepo provides an interface for interacting with a storage solution
type MessageRepo interface {
	Fetch(ID string) []*Message
	FetchX(ID string, x int) []*Message
	Put(ID string, msgs ...*Message) error
	Delete(ID string) error
}

// Backupper will backup messages to its store after
// A: the set time interval has passed or
// B: the in memeory cache has reached its capacity
//
// NOTE: id is used as the key that the data will saved under
// in the data store.
type Backupper struct {
	cache    []*Message
	store    MessageRepo
	ticker   time.Ticker
	interval int
	id       string
	busy     bool
}

// Run starts the backupper and listens forever on its ticker channel,
// calling backup at the desired interval.
// When run is terminated using context, we check if a backup is already in progess
// and if not we backup before terminating.
func (b *Backupper) Run(ctx context.Context) {
loop:
	for {
		select {
		case <-b.ticker.C:
			b.busy = true

			// TODO: handle error from backup
			b.Backup()

			b.busy = false
		case <-ctx.Done():
			if !b.busy {
				b.Backup()
			}

			b.ticker.Stop()

			b.cache = nil // free any leftover memory since we reuse the cache

			break loop
		}
	}
}

// Hold stores any number of messages inside its in mem cache.
func (b *Backupper) Hold(msgs ...*Message) {
	b.cache = append(b.cache, msgs...)
}

// Backup purges all messages from cache into store
// then reuses the slice. All memory is eventally freed,
// when the run method ends.
func (b *Backupper) Backup() error {
	err := b.store.Put(b.id, b.cache...)

	if err != nil {
		return err
	}

	// reuse the cache
	// NOTE: see here about possibly mem-leaks with this method
	// https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
	b.cache = b.cache[:0]

	return nil
}
