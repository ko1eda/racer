package racer

import (
	"net/http"
)

// Subscriber represents a client who is elegable to recieve broadcasts from the broker
// A client must be able to register with the broker and provide a buffered channel which the broker
// can return information on.
type Subscriber interface {
	Register(broadcast chan<- []byte, unregister chan chan<- []byte) (send chan<- []byte)
}

// Client a client must be able to read from some kind of connection, whether it be tcp, rcp, webscoket etc
// It must also run a process that is able to dameonize well it reads to and write from said connection
type Client interface {
	ReadFromCon()
	WriteToCon()
	Run(w http.ResponseWriter, r *http.Request)
}

// BrokeredClient is a client that utilizes the message broker to send and recieve updates.
// A client that implements this interface must register with a running broker before it itself can be run.
type BrokeredClient interface {
	Client
	Subscriber
}

// type BrokerManager struct {
// 	map[string]*Broker
// }
