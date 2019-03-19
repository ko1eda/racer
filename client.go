package racer

// import (
// 	"math/rand"

// 	"github.com/tinylttl/racer/broker"
// )

// // Connection is the source of data to and from the client and server.
// // The default connection type for racer is socket.
// type Connector interface {
// 	Read() <-chan *Message
// 	Write() chan<- *Message
// }

// // Broadcaster can broadcast messages to other listening client goroutines
// type Broadcaster interface {
//	Register()
// 	Broadcast() chan<- *broker.Message
// 	Unregister() chan chan<- *broker.Message
// }

// // Message is data that is sent as json through the connection.
// type Message struct {
// 	Timestamp int64  `json:"timestamp"`
// 	Sent      string `json:"sent"`
// 	Body      string `json:"body"`
// 	SenderID  int    `json:"senderID"`
// }

// type Client struct {
// 	Broadcaster
// 	Connector
// 	// MessageRepo
// }

// // NewClient returns a new Chat client instance
// func NewClient(b Broadcaster, conn Connection) *Chat {
// 	return &Chat{id: rand.Intn(100000), send: make(chan *broker.Message, 1), Broadcaster: b, Connection: conn}
// }

// Run starts a deamonized chat instance
// func (c *client) Run(w http.ResponseWriter, r *http.Request) {
// 	// Hijack the connection from the server so the client can communicate directly through the socket
// 	c.Hijack(w http.ResponseWriter, r *http.Request)

// 	go func() {
// 		for msg := range c.Read() {
// 			c.Broadcast() <- &broker.Message{Payload: msg}
// 		}
// 		// shutdown the client because the connection was closed
// 	}()
// 	go func() {
// 		for msg := range c.send {
// 			c.Write() <- msg
// 		}
// 	}()
// loop:
// 	for {
// 		select {
// 		case cmsg, ok <- c.Read() :
// 			if !ok {
// 				c.Unregister() <- c.send
// 				break loop
// 			}
// 			c.Broadcast() <- &broker.Message{Payload: cmsg}
// 		case bmsg, ok <- c.send :
// 			if !ok {
// 				break loop
// 			}
// 			c.Write() <-bmsg.Payload.(*Message)
// 		}
// 	}
// }
