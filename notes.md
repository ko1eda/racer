### Design


### What i know
+ Each client gets its own socket, which it is constantly reading from for new sends+
+ When a send occurs, it broadcasts that send on a Room broadcast channel
+ The room is responsible for taking all messages on a given clients socket and broadcasting them
+ to all connected sockets


### Todo
+ Write tests to veryify the clients are being created correctly in the correct rooms(hubs/managers/broker)
+ clean up rocketd directory and write basic route matching with chi
+ transmit json data through the socket
+ Make Room initialization happen inside a handler function
+ Add database support for retreival of chats for a given id 
+ research concurrent design patterns 
+ research backing up data on client interruption (how can I block for a certain period of time and wait for the client to reconnect and if not store their data )