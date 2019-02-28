package racer

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Client a client must be able to read from some kind of connection, whether it be tcp, rcp, webscoket etc
// It must also run a process that is able to dameonize well it reads to and write from said connection
type Client interface {
	ReadFromCon()
	WriteToCon()
	Run(w http.ResponseWriter, r *http.Request)
}

// BrokeredClient is a client that utilizes the message broker to send and recieve updates
type BrokeredClient interface {
	Client
	Subscriber
}

// A Chat represents a wrapper for a socket, its job is to read and write to the socket
// It reads any data from the socket and then broadcasts it onto a brokers global broadcast channel
// the broker then sends the message to all clients on their individual send buses
// the clients then push the messages into their sockets one by one propegating the message from one client to all other sockets
// this works because sockets are like channels they are buffered, when one is being written to it is blocked, and when it is read from the data inside the
// socket is removed.. So we must rebroadcast the data  when we read it to all listneing parties who will read it and then write it from their send channels
// https://stackoverflow.com/questions/14241235/what-happens-when-i-write-data-to-a-blocking-socket-faster-than-the-other-side
type Chat struct {
	con        *websocket.Conn // a con is used to read from and write to which intern updates all clients in the broker
	broadcast  chan<- []byte
	unregister chan chan<- []byte
	send       chan []byte // each client has their own unique send channel for sending data from the broadcast channel into the con
	id         int
}

// A Message represents chat data sent between users in a broker
type Message struct {
	// sent        time.Time
	// retrieved   time.Time
	// senderID    int
	// retrieverID int
	body string
}

// NewChat returns a new Chat client instance
func NewChat() BrokeredClient {
	return &Chat{id: rand.Intn(100000), send: make(chan []byte, 256)}
}

// Register a chat client with a broker, all chats on the brokers channel will recieve the same updates
func (c *Chat) Register(broadcast chan<- []byte, unregister chan chan<- []byte) chan<- []byte {
	c.broadcast = broadcast
	c.unregister = unregister
	return c.send
}

// Run starts a deamonized chat instance
func (c *Chat) Run(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	c.con = con

	go c.ReadFromCon()
	go c.WriteToCon()
}

// ReadFromCon the client reads from its connection and sends the message to any other sibling clients through its brokers broadcast channel
func (c *Chat) ReadFromCon() {
	// Defere the closing of the con and deregistration to when this function terminates
	// it will only terminate if the client disconnects or there is an error
	defer func() {
		c.con.Close()
		c.unregister <- c.send
	}()

	// The maximum bytes our read routines can read in from the con is 512 bytes so 512 1 byte asci characters
	c.con.SetReadLimit(maxMessageSize)
	c.con.SetReadDeadline(time.Now().Add(pongWait))

	// Every time a pong occurs on the con, our read routine will add more time before it times out
	c.con.SetPongHandler(func(string) error { c.con.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.con.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(c.id, ": ", string(message))
		c.broadcast <- message
	}
}

// WriteToCon the client should use data from their send channel to update their con
func (c *Chat) WriteToCon() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.con.Close()
	}()

	// if their are no messages from any other clients the ticker pings all other members of the con
	// triggering their pong handlers which intern rerefreshes their read deadlines
	for {
		select {
		case msg, ok := <-c.send:
			fmt.Println("length of send channel inside write loop ", len(c.send))
			c.con.SetWriteDeadline(time.Now().Add(writeWait))

			// If the clients send channel has been closed by the broker then there was an error
			// and this peer will send a close message to the con, meaining that it (the clients) connection will be closed
			if !ok {
				c.con.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("The clients send channel was closed ")
				return
			}

			// Get the next io writer from the con, only one writer can be held by any given function
			// TODO: Change to binary protocol and use json
			w, err := c.con.NextWriter(websocket.TextMessage)

			if err != nil {
				c.con.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// write the message from the current clients send channel (which is filled with data from the brokers broadcast channel) to the con
			w.Write(msg)

			// Check to see that there are no built up messages in the send channel
			// Send is a Buffered channel so it is possible that there will be more bytes on the channel that have built up
			// after this select happend
			for i := 0; i < len(c.send); i++ {
				w.Write([]byte(`\n`))
				w.Write(<-c.send)
			}

			// When we are finished writing to the con
			// close the writer that we were using
			if err := w.Close(); err != nil {
				log.Println("The writer could not be closed")
				return
			}
		case <-ticker.C:
			c.con.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.con.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
