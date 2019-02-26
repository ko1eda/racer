/*

	Package Rocket is a

	Rooms store multiple client connections and
	broadcast to them if there are any updates on their broadcast channels.
	A room is like a client manager.


	Client a client represents an individual connection to a room.
	Each client upgrades the incomming connection to a ws connection and spawns 2 go routines
	1 for reading from the socket and broadcasting to its rooms broadcast channel when data is found in the socket
	1 for writing from the clients send channel into the socket

	The clients send channel is populated by its Rooms broadcast channel when the client reads a message from the socket (which was sent via a write to the socket from the js client)

*/

package rocket
