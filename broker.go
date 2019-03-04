package racer

import (
	"fmt"
	"sync"
)

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
loop:
	for {
		select {
		case client := <-b.register:
			b.subscribers[client] = true

		case unregistered := <-b.unregister:
			delete(b.subscribers, unregistered)
			close(unregistered)

			if len(b.subscribers) == 0 {
				break loop
			}

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

// BrokerManager keeps a mapping of chatIDs and brokers
// it ensures that only one broker may be active for a given chatID
type BrokerManager struct {
	brokerm map[string]*Broker
	mu      sync.Mutex
}

// NewManager creates a new BrokerManager
func NewManager() *BrokerManager {
	return &BrokerManager{brokerm: make(map[string]*Broker)}
}

// // Register registers a new broker with the manager, it returns true
// // if the registration succeeded and false if not
// func (bm *BrokerManager) Register(key string, b *Broker) bool {
// 	if _, exists := bm.BrokerExists(key); !exists {
// 		bm.mu.Lock()
// 		fmt.Println("called register")
// 		bm.brokerm[key] = b
// 		bm.mu.Unlock()
// 		return true
// 	}

// 	return false
// }

// Unregister unregisters a new broker with the manager, deleting its key
// from its map it returns true if the key was found and delete false if it was not found
func (bm *BrokerManager) Unregister(key string) bool {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if _, exists := bm.brokerm[key]; exists {
		delete(bm.brokerm, key)
		return true
	}

	return false
}

// BrokerExists uses a Read lock to check if a broker already exists for a given key
// It returns a boolean true if it does or false if does not and closes the readlock.
func (bm *BrokerManager) BrokerExists(key string) (*Broker, bool) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	broker, exists := bm.brokerm[key]
	fmt.Println("called brokerexists ", exists)

	if exists {
		return broker, true
	}

	return nil, false
}

// ExistsOrNew checks to see if a broker exists for a given chatID.
// if it does not exist already, it creates it and adds an entry to its map and returns a new broker and false.
// If it does already exist it returns the entry from the map and true
//
// It uses a sync.Mutex to keep its underlying map reads and writes thread safe.
func (bm *BrokerManager) ExistsOrNew(key string) (*Broker, bool) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	broker, exists := bm.brokerm[key]

	if exists {
		fmt.Println("BROKER EXISTS SO WE RETURN IT")
		return broker, true
	}

	broker = NewBroker()
	bm.brokerm[key] = broker

	fmt.Println("MAKE NEW BROKER")

	return broker, false
}
