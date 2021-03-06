package broker

import (
	"time"
)

// A Topic represents a connection hub, anything registered with a topic will recieve updates
// every time a message is pushed to its broadcast channel. A topic must be started in order for it
// to register subscribers and broadcast messages.
type Topic struct {
	subscribers map[chan<- *Message]bool
	register    chan chan<- *Message
	broadcast   chan *Message
	unregister  chan chan<- *Message
	ID          string
}

// Message is sent through the brokers broadcast channel and relayed to any listeners through
// their respective send channels.
type Message struct {
	Recieved time.Time   // when the topic got the message on its broadcast channel
	Sent     time.Time   // when the topic sent the message to all its registered client channels
	Payload  interface{} // any clients using the same topic should be expecting the same type of message
}

// NewTopic creates a new Topic
func NewTopic(ID string) *Topic {
	return &Topic{
		ID:          ID,
		subscribers: make(map[chan<- *Message]bool),
		broadcast:   make(chan *Message),
		register:    make(chan chan<- *Message),
		unregister:  make(chan chan<- *Message),
	}
}

// Register registers a new send channel with the topic. Clients will recieve on this channel
// whenever there is a message sent to the brokers broadcast channel.
func (t *Topic) Register() chan chan<- *Message { return t.register }

// Broadcast exposes a topics internal broadcast channel.
// Use this to send messages to other clients that subscribe to this topic.
func (t *Topic) Broadcast() chan<- *Message { return t.broadcast }

// Unregister exposes a topics internal channel for unregistering clients.
func (t *Topic) Unregister() chan chan<- *Message { return t.unregister }

// Start starts the Topic in a blocking state, it will listen on all its channels
// and select an action based on the currently active channel.
// If a new client is registered to the Topic, it will update its map of subscribers
// If a client is unregistered from the Topic it will remove it from its list of subscribers and close its channel
// If the Brokers boradcast channel recieves a message, it will relay that message to all subscribers in its map through their respective send channels
func (t *Topic) Start() {
loop:
	for {
		select {
		case client := <-t.register:
			t.subscribers[client] = true

		case unregistered := <-t.unregister:
			delete(t.subscribers, unregistered)
			close(unregistered)

			if len(t.subscribers) == 0 {
				break loop
			}

		case msg := <-t.broadcast:
			for client := range t.subscribers {
				select {
				case client <- msg:
					// msg.Recieved = time.Now() // does cause a race condition
				default:
					close(client)
					delete(t.subscribers, client)
				}
			}
		}
	}
}
