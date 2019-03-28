package broker

import (
	"sync"

	"github.com/tinylttl/racer/id"
)

// Broker keeps a mapping of chatIDs and brokers
// it ensures that only one topic may be active for a given chatID
type Broker struct {
	topics map[string]*Topic
	//idgen	racer.id
	mu     sync.Mutex
}

// NewBroker creates a new Broker. A new map is intialized by default if WithMap option is not passed in.
func NewBroker(opts ...func(*Broker)) *Broker {
	b := Broker{topics: make(map[string]*Topic)}

	for _, opt := range opts {
		opt(&b)
	}

	return &b
}

// WithMap allows you to pass in your own map that the manager will use to map keys to active brokers
// this can be useful in testing where you would like direct access to the managers internal mappings
func WithMap(m map[string]*Topic) func(*Broker) {
	return func(b *Broker) {
		b.topics = m
	}
}

// NewTopic returns a newly initialized topic with a unique identifier. It also starts the topic. This is a convienience method for NewTopic()
func (b *Broker) NewTopic() *Topic {
	g, _ := id.NewGenerator() // this should be injected or be a part of the broker struct
	id, _ := g.NewID()
	t := NewTopic(id)

	go func() {
		t.Start()
	}()

	b.Add(id, t)

	return t
}

// Add adds a new topic to the brokers map of active topics.
// Returns true if it could be added, false if there was already a topic with that key.
func (b *Broker) Add(key string, t *Topic) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, exists := b.topics[key]
	if !exists {
		b.topics[key] = t
		return true
	}

	return false
}

// Lookup searches its map for a running topic with the given id.
// The callback will be called regardless of whether or not a topic is found.
//
// If a topic is found, we will pass it in and set found to true.
// If a topic is not found a new one is created and registered in the map, found will then be set to false and the new topic is passed to the cb.
//
// NOTE: If you would like to remove a topic from the manager, make sure you always call the BrokerManagers Remove method as it is thread safe.
func (b *Broker) Lookup(key string, cb func(found bool, b *Topic)) {
	b.mu.Lock()
	topic, exists := b.topics[key]

	if !exists {
		topic := NewTopic(key)
		b.topics[key] = topic
		b.mu.Unlock()

		cb(false, topic)
		return
	}

	b.mu.Unlock()
	cb(true, topic)
}

// Size returns the size of the brokers underlying map
func (b *Broker) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.topics)
}

// Remove removes a topic from the manager deleting the key from its map
// It returns true if the key was found and deleted false if it was not found.
func (b *Broker) Remove(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.topics[key]; exists {
		delete(b.topics, key)
		return true
	}

	return false
}

// Exists uses a lock to check if a topic already exists for a given key
// It returns a boolean true if it does or false if does not and closes the lock.
func (b *Broker) Exists(key string) (*Topic, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	t, exists := b.topics[key]

	if exists {
		return t, true
	}

	return nil, false
}
