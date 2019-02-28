package racer

// Subscriber represents a client who is elegable to recieve broadcasts from the broker
// A client must be able to register with the broker and provide a buffered channel which the broker
// can return information on.
type Subscriber interface {
	Register(broadcast chan<- []byte, unregister chan chan<- []byte) (send chan<- []byte)
}

// A Broker represents a connection hub, anything registered with a broker will recieve updates
// every time a message is pushed to its broadcast channel
type Broker struct {
	id          string
	subscribers map[chan<- []byte]bool
	broadcast   chan []byte
	register    chan chan<- []byte
	unregister  chan chan<- []byte
}

// NewBroker creates a new Broker
func NewBroker() *Broker {
	return &Broker{
		id:          "MYCOOLBroker",
		subscribers: make(map[chan<- []byte]bool),
		broadcast:   make(chan []byte),
		register:    make(chan chan<- []byte),
		unregister:  make(chan chan<- []byte),
	}
}

// Start starts the Broker in a blocking state, it will listen on all its channels
// and select an action based on the currently active channel.
// If a new client is registered to the Broker, it will update its map of subscribers
// If a client is unregistered from the Broker it will remove it from its list of subscribers and close its channel
// If the Brokers boradcast channel recieves a message, it will relay that message to all subscribers in its map through their respective send channels
func (b *Broker) Start() {
	for {
		select {
		case client := <-b.register:
			b.subscribers[client] = true

		case unregistered := <-b.unregister:
			delete(b.subscribers, unregistered)
			close(unregistered)

		case msg := <-b.broadcast:
			for client := range b.subscribers {
				select {
				case client <- msg:
				default:
					close(client)
					delete(b.subscribers, client)
				}
			}
		}
	}
}

// RegisterSubscriber brokers a client connection - that is it adds a client to its list of subscribers
// to broadcast to. And upgrades the client to a socket connection
func (b *Broker) RegisterSubscriber(s Subscriber) { b.register <- s.Register(b.broadcast, b.unregister) }
