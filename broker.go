package racer

import (
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
func NewBroker(id string) *Broker {
	return &Broker{
		id:          id,
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

// NewManager creates a new BrokerManager. A new map is intialized by default if WithMap option is not passed in.
func NewManager(opts ...func(*BrokerManager)) *BrokerManager {
	bm := BrokerManager{
		brokerm: make(map[string]*Broker),
	}

	for _, opt := range opts {
		opt(&bm)
	}

	return &bm
}

// WithMap allows you to pass in your own map that the manager will use to map keys to active brokers
// this can be useful in testing where you would like direct access to the managers internal mappings
func WithMap(m map[string]*Broker) func(*BrokerManager) {
	return func(bm *BrokerManager) {
		bm.brokerm = m
	}
}

// Lookup searches its map for a running broker with the given id.
// The callback will be called regardless of whether or not a broker is found.
//
// If a broker is found, we will pass it in and set found to true.
// If a broker is not found a new one is created and registered in the map, found will then be set to false and the new broker is passed to the cb.
//
// NOTE: If you would like to remove a broker from the manager, make sure you always call the BrokerManagers Remove method as it is thread safe.
func (bm *BrokerManager) Lookup(key string, cb func(found bool, b *Broker)) {
	bm.mu.Lock()
	broker, exists := bm.brokerm[key]

	if !exists {
		broker := NewBroker(key)
		bm.brokerm[key] = broker
		bm.mu.Unlock()

		cb(false, broker)
		return
	}

	bm.mu.Unlock()
	cb(true, broker)
}

// Remove removes a broker from the manager deleting the key from its map
// It returns true if the key was found and deleted false if it was not found.
func (bm *BrokerManager) Remove(key string) bool {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if _, exists := bm.brokerm[key]; exists {
		delete(bm.brokerm, key)
		return true
	}

	return false
}

// Exists uses a lock to check if a broker already exists for a given key
// It returns a boolean true if it does or false if does not and closes the lock.
func (bm *BrokerManager) Exists(key string) (*Broker, bool) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	b, exists := bm.brokerm[key]

	if exists {
		return b, true
	}

	return nil, false
}

// Result contains the result of a lookup from the broker manager
// It will contain a Broker and a boolean field Found which will be true
// if the broker was found and false if it was not.
// If the broker was not found the manager creates one and stores it under the lookup key.
// type Result struct {
// 	*Broker
// 	Found bool
// }
// // Lookup exposes an unbuffered channel on the broker manager
// // it takes a string key and will check the managers lookup table
// // to determine if a broker exists for the given key
// func (bm *BrokerManager) Lookup() chan<- string {
// 	go bm.runLookup()
// 	return bm.lookup
// }

// // Result returns a channel that containts the results from each lookup
// func (bm *BrokerManager) Result() <-chan Result {
// 	return bm.result
// }

// func (bm *BrokerManager) runLookup() {
// 	for {
// 		select {
// 		case key := <-bm.lookup:
// 			broker, exists := bm.brokerm[key]
// 			if !exists {
// 				broker = key
// 				bm.brokerm[key] = broker

// 				bm.result <- Result{broker, false}
// 				break
// 			}

// 			bm.result <- Result{broker, true}
// 		}
// 	}
// }
// Lookup exposes an unbuffered channel on the broker manager
// it takes a string key and will check the managers lookup table
// to determine if a broker exists for the given key

// // Result returns a channel that containts the results from each lookup
// func (bm *BrokerManager) Result() <-chan Result {
// 	return bm.result
// }

// func (bm *BrokerManager) runLookup() {
// 	for {
// 		select {
// 		case key := <-bm.lookup:
// 			broker, exists := bm.brokerm[key]
// 			if !exists {
// 				broker = key
// 				bm.brokerm[key] = broker

// 				bm.result <- Result{broker, false}
// 				break
// 			}

// 			bm.result <- Result{broker, true}
// 		}
// 	}
// }
