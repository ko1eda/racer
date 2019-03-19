package gorilla

// // ReadFromCon the client reads from its connection and sends the message to any other sibling clients through its brokers broadcast channel
// func (c *Chat) readFromCon() {
// 	// Defere the closing of the con and deregistration to when this function terminates
// 	// it will only terminate if the client disconnects or there is an error
// 	defer func() {
// 		c.con.Close()
// 		c.Unregister() <- c.send
// 	}()

// 	// The maximum bytes our read routines can read in from the con is 512 bytes so 512 1 byte asci characters
// 	// Every time a pong occurs on the con, our read routine will add more time before it times out
// 	c.con.SetReadLimit(maxMessageSize)
// 	c.con.SetReadDeadline(time.Now().Add(pongWait))
// 	c.con.SetPongHandler(func(string) error { c.con.SetReadDeadline(time.Now().Add(pongWait)); return nil })

// 	// remember that this *message has a val of nil, so to use unmarshall
// 	// we need to pass its refrence, so that the mem address *message points to can be filled with a value
// 	// if we just pass chatmsg it would be passing nil to the unmarshall func.
// 	// we could also declare message as a concrete type (without *) and pass its &refrence, doing so intializes the 0 val for chatmsg
// 	// and we pass it the address of that.
// 	brokermsg := broker.Message{Sent: time.Now()}
// 	chatmsg := message{Sent: time.Now().Format("01/02/06 15:04 pm"), Timestamp: time.Now().UTC().Unix()}
// 	for {
// 		err := c.con.ReadJSON(&chatmsg)

// 		// v := json.Unmarshal()
// 		// fmt.Printf("Chat message %#v\n", chatmsg)
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		brokermsg.Payload = &chatmsg
// 		// fmt.Printf("Chat message %#v\n", brokermsg)
// 		c.Broadcast() <- &brokermsg
// 	}
// }

// // WriteToCon the client should use data from their send channel to update their con
// func (c *Chat) writeToCon() {
// 	ticker := time.NewTicker(pingPeriod)

// 	defer func() {
// 		ticker.Stop()
// 		c.con.Close()
// 	}()

// 	// if their are no messages from any other clients the ticker pings all other members of the con
// 	// triggering their pong handlers which intern rerefreshes their read deadlines
// 	// var chatmsg *message
// 	for {
// 		select {
// 		case brokermsg, ok := <-c.send:
// 			c.con.SetWriteDeadline(time.Now().Add(writeWait))

// 			// If the clients send channel has been closed by the broker then there was an error
// 			// and this peer will send a close message to the con, meaining that it (the clients) connection will be closed
// 			if !ok {
// 				c.con.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			m, ok := brokermsg.Payload.(*message)

// 			if !ok {
// 				c.con.WriteMessage(websocket.CloseMessage, []byte{})
// 				log.Println("Error payload was unexpected type")
// 				return
// 			}

// 			err := c.con.WriteJSON(m)
// 			if err != nil {
// 				c.con.WriteMessage(websocket.CloseMessage, []byte{})
// 				log.Println("Error writing json to conn. ", err)
// 				return
// 			}

// 			// Check to see that there are no built up messages in the send channel
// 			// Send is a Buffered channel so it is possible that there will be more bytes on the channel that have built up
// 			// after this select happend
// 			// for i := 0; i < len(c.send); i++ {

// 			// }
// 		case <-ticker.C:
// 			c.con.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.con.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
