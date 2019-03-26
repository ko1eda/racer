package gorilla

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tinylttl/racer"
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

	// time format for a racer.Message
	timeFmt = "01/02/06 3:04 pm"
)

// Check that the Connector implementation can be assigned to a racer.Connector interface
// this ensures they are compatable
// var _ racer.Connector = (*Connector)(nil) // written this way I think yields no extra allocations but i'm not really sure exactly whats happening so I switched to the way I understand
var _ racer.Connector = &Connector{}

// Connector represents a single socket connection that can be held by a client
// It provides read and write channels that can be used to read data from the socket
// and write data to the socket, respectively.
type Connector struct {
	conn         *websocket.Conn
	rchan, wchan chan *racer.Message // read and write channels for communicating messages recieved through socket
}

// NewConnection returns a connector with a newly upgraded socket connection
// Connector implements racer.Connector interface
func NewConnection(w http.ResponseWriter, r *http.Request) *Connector {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	return &Connector{conn: conn, rchan: make(chan *racer.Message, 10), wchan: make(chan *racer.Message, 10)}
}

// Read reads data from a socket and returns it on a read-only channel
func (c *Connector) Read() <-chan *racer.Message {
	go func() {
		// Defer the closing of the con and deregistration to when this function terminates
		// it will only terminate if the client disconnects or there is an error
		defer func() {
			c.conn.Close()
			close(c.rchan)
		}()

		// The maximum bytes our read routines can read in from the con is 512 bytes so 512 1 byte asci characters
		// Every time a pong occurs on the con, our read routine will add more time before it times out
		c.conn.SetReadLimit(maxMessageSize)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		// remember that this *message has a val of nil, so to use unmarshall
		// we need to pass its refrence, so that the mem address *message points to can be filled with a value
		// if we just pass chatmsg it would be passing nil to the unmarshall func.
		// we could also declare message as a concrete type (without *) and pass its &refrence, doing so intializes the 0 val for chatmsg
		// and we pass it the address of that.
		chatmsg := racer.Message{Sent: time.Now().Format(timeFmt), Timestamp: time.Now().UTC().Unix()}
		for {
			err := c.conn.ReadJSON(&chatmsg)
			// v := json.Unmarshal()
			// fmt.Printf("Chat message %#v\n", chatmsg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}
			c.rchan <- &chatmsg
			// fmt.Printf("Chat message %#v\n", brokermsg)
		}
	}()

	return c.rchan
}

// Write the client should use data from their send channel to update their con
func (c *Connector) Write() chan<- *racer.Message {
	go func() {
		ticker := time.NewTicker(pingPeriod)

		defer func() {
			ticker.Stop()
			c.conn.Close()
		}()

		// if their are no messages from any other clients the ticker pings all other members of the con
		// triggering their pong handlers which intern rerefreshes their read deadlines
		// var chatmsg *message
		for {
			select {
			case msg, ok := <-c.wchan:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))

				// If the clients send channel has been closed by the broker then there was an error
				// and this peer will send a close message to the con, meaining that it (the clients) connection will be closed
				if !ok {
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				err := c.conn.WriteJSON(msg)
				if err != nil {
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					log.Println("Error writing json to conn. ", err)
					return
				}

				// Check to see that there are no built up messages in the send channel
				// Send is a Buffered channel so it is possible that there will be more bytes on the channel that have built up
				// after this select happend
				// for i := 0; i < len(c.send); i++ {

				// }
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	return c.wchan
}
