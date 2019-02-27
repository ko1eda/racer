package rocket

// A Room represents an inidividual socket connection that can be shared
// between multiple authenticated users. Authentication happens through our authentication server not part of this application
// The room will interact with
type Room struct {
	id         string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	repo       string // msgrepo that interacts with the database
}

// NewRoom creates a new room
func NewRoom() *Room {
	return &Room{
		id:         "MYCOOLROOM",
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// StartBroadcast starts the room in a blocking state, it will listen on all its channels
// and select an action based on the currently active channel.
// If a new client is registered to the room, it will update its map of clients
// If a client is unregistered from the room it will remove it from its list of clients and close its buffered send channel
// If the rooms boradcast channel recieves a message, it will relay that message to all clients in its map through their respective send channels
func (r *Room) StartBroadcast() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true

		case unregistered := <-r.unregister:
			delete(r.clients, unregistered)
			close(unregistered.send)

		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
