package racer

import (
	"fmt"
	"math/rand"

	"github.com/tinylttl/racer/broker"
)

// Client represents that chat client. Everytime a new connection is made
// to the server, a new client is created.
type Client struct {
	Broadcaster Broadcaster
	Conn        Connector
	ID          string
	Receive     chan *broker.Message // Used to Receive messages from the broadcaster
	//MsgRepo 	Repo interface
}

// Message is data that is sent as json through the connection.
type Message struct {
	Timestamp int64  `json:"timestamp"`
	Sent      string `json:"sent"`
	Body      string `json:"body"`
	SenderID  int    `json:"senderID"`
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

// Run starts two goroutines, one to read incoming messages from the Clients connection
// the other is to read messages recieved from the clients broadcaster, and send them through
// the clients connection
func (c *Client) Run() {
	go func() {
		for msg := range c.Conn.Read() {
			c.Broadcaster.Broadcast() <- &broker.Message{Payload: msg}
		}
		c.Broadcaster.Unregister() <- c.Receive // shutdown the client because the connection was closed
	}()

	go func() {
		for bmsg := range c.Receive {
			c.Conn.Write() <- bmsg.Payload.(*Message)
		}
	}()
}
