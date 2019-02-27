package rocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// A Client represents a wrapper for a socket, its job is to read and write to the socket
// It reads any data from the socket and then broadcasts it onto a rooms global broadcast channel
// the room then sends the message to all clients on their individual send buses
// the clients then push the messages into their sockets one by one propegating the message from one client to all other sockets
// this works because sockets are like channels they are buffered, when one is being written to it is blocked, and when it is read from the data inside the
// socket is removed.. So we must rebroadcast the data  when we read it to all listneing parties who will read it and then write it from their send channels
// https://stackoverflow.com/questions/14241235/what-happens-when-i-write-data-to-a-blocking-socket-faster-than-the-other-side
type Client struct {
	socket *websocket.Conn // a socket is used to read from and write to which intern updates all clients in the room
	room   *Room
	send   chan []byte // each client has their own unique send channel for sending data from the broadcast channel into the socket
	id     int
}

// A Message represents chat data sent between users in a room
type Message struct {
	// sent        time.Time
	// retrieved   time.Time
	// senderID    int
	// retrieverID int
	body string
}

func (c *Client) readFromSocket() {
	// When this function finishes its execution
	// close the socket if it is open and unregister this client from its Publisher aka its room
	defer func() {
		c.socket.Close()
		c.room.unregister <- c
	}()

	// The maximum bytes our read routines can read in from the socket is 512 bytes so 512 1 byte asci characters
	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))

	// Every time a pong occurs on the socket, our read routine will add more time before it times out
	c.socket.SetPongHandler(func(string) error { c.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(c.id, ": ", string(message))
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.room.broadcast <- message
	}
}

func (c *Client) writeToSocket() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()

	// if their are no messages from any other clients the ticker pings all other members of the socket
	// triggering their pong handlers which intern rerefreshes their read deadlines
	for {
		select {
		case msg, ok := <-c.send:
			fmt.Println("length of send channel inside write loop ", len(c.send))
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))

			// If the clients send channel has been closed by the room then there was an error
			// and this peer will send a close message to the socket, meaining that it (the clients) connection will be closed
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("The clients send channel was closed ")
				return
			}

			// Get the next io writer from the socket, only one writer can be held by any given function
			// TODO: Change to binary protocol and use json
			w, err := c.socket.NextWriter(websocket.TextMessage)

			if err != nil {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// write the message from the current clients send channel (which is filled with data from the rooms broadcast channel) to the socket
			w.Write(msg)

			// Check to see that there are no built up messages in the send channel
			// Send is a Buffered channel so it is possible that there will be more bytes on the channel that have built up
			// after this select happend
			for i := 0; i < len(c.send); i++ {
				w.Write([]byte(`\n`))
				w.Write(<-c.send)
			}

			// When we are finished writing to the socket
			// close the writer that we were using
			if err := w.Close(); err != nil {
				log.Println("The writer could not be closed")
				return
			}
		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

// StartClientInRoom spawns new client read and write processes and sets its room to
// the provided room struct
func StartClientInRoom(id int, room *Room, w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	c := Client{socket, room, make(chan []byte, 256), id}

	c.room.register <- &c

	go c.readFromSocket()
	go c.writeToSocket()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
