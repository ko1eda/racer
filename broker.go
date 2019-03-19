package racer

import "github.com/tinylttl/racer/broker"

// Broadcaster can broadcast messages to other listening client goroutines
type Broadcaster interface {
	Register() chan chan<- *broker.Message
	Broadcast() chan<- *broker.Message
	Unregister() chan chan<- *broker.Message
}
