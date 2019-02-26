package rocket

import "time"

// type Publisher struct {
// 	socket string
// }

// A Room represents an inidividual socket connection that can be shared
// between multiple authenticated users. Authentication happens through our authentication server not part of this application
// The room will interact with
type Room struct {
	id         string
	clients    map[*Client]bool
	broadcast  chan []Message
	register   chan *Client
	unregister chan *Client
	repo       string // msgrepo that interacts with the database
}

// A Client represents a wrapper for a socket, its job is to read from the socket sending any data it finds to its Rooms broadcast channel
// Any data that it receives from the boradcast channel is to be written into the socket which intern is pushed to the other clients in the chat
type Client struct {
	send   chan []Message // each client has their own unique send channel for sending data from the broadcast channel into the socket
	socket string         // a socket is used to read from and write to which intern updates all clients in the room
}

// // DB holds a connection to the database where we store the rooms and message data
// type DB struct {
// 	db string // bolt db
// }

// MsgRepo interacts with our database and performs the relevent actions like fetchall and store
type MsgRepo struct {
	db string // badget db
}

// A Message represents chat data sent between users in a room
type Message struct {
	sent        time.Time
	retrieved   time.Time
	senderID    int
	retrieverID int
	body        string
}
