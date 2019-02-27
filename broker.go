package racer

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// A Broker represents an inidividual socket connection that can be shared
// between multiple authenticated users. Authentication happens through our authentication server not part of this application
// The room will interact with
type Broker struct {
	id         string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	repo       string // msgrepo that interacts with the database
}

// NewBroker creates a new Broker
func NewBroker() *Broker {
	return &Broker{
		id:         "MYCOOLBroker",
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// StartBroadcast starts the Broker in a blocking state, it will listen on all its channels
// and select an action based on the currently active channel.
// If a new client is registered to the Broker, it will update its map of clients
// If a client is unregistered from the Broker it will remove it from its list of clients and close its buffered send channel
// If the Brokers boradcast channel recieves a message, it will relay that message to all clients in its map through their respective send channels
func (b *Broker) StartBroadcast() {
	for {
		select {
		case client := <-b.register:
			b.clients[client] = true

		case unregistered := <-b.unregister:
			delete(b.clients, unregistered)
			close(unregistered.send)

		case msg := <-b.broadcast:
			for client := range b.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(b.clients, client)
				}
			}
		}
	}
}

// BrokerClientCon brokers a client connection - that is it adds a client to its list of clients
// to broadcast to. And upgrades the client to a socket connection
func (b *Broker) BrokerClientCon(c *Client, w http.ResponseWriter, r *http.Request, options ...func(*Client)) {
	socket, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	// replace this with a uuid
	rand.Seed(time.Now().UnixNano())
	rand.Intn(100000)
	c = &Client{socket, b, make(chan []byte, 256), rand.Intn(100000)}

	c.broker.register <- c

	// Extra configuration for the client broker
	for _, fn := range options {
		fn(c)
	}

	go c.ReadFromSocket()
	go c.WriteToSocket()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// type Client interface {
// 	NewClient(broadcast ->chan []byte, register ->chan *Client, unregister ->chan *Client, id int, socket )
// 	ReadFromSocket()
// 	WriteToSocket()
// 	GetSendChannel() * -> chan []byte
// }
