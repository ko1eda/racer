### Design


### Setting errors on response writier
+ https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter

### What i know
+ Each client gets its own socket, which it is constantly reading from for new sends+
+ When a send occurs, it broadcasts that send on a Room broadcast channel
+ The room is responsible for taking all messages on a given clients socket and broadcasting them
+ to all connected sockets


### Todo
+ Write tests to veryify the clients are being created correctly in the correct rooms(hubs/managers/broker)
    + Test1: Test that two requests to the same chatid creates only one new entry in our MAP
    + Test2: Check that two concurrent requests to the same chatID creates only one new entry in our MAP
    + Test3: Check that two concurrent exits from the same chat Does not cause a race condition (writing to brokermap)

- Stop the broker if all clients have left and remove it from the map so it doesn't stay running for no reason (needs testing)

+ Add os signal package to our racerd and use context with each route, this way we can gracefully shutdown the program if something happens 
    + Put server on a select with a signal channel
    + If the server fails cancel parent context, if the signal is canceled cancel parent context 





+ Work on implementing client interface to decouple brokers from clients 
+ Figure out what time variables are doing in client.go




+ transmit json data through the socket
+ Make Room initialization happen inside a handler function
+ Add database support for retreival of chats for a given id 
+ research concurrent design patterns 
+ research backing up data on client interruption (how can I block for a certain period of time and wait for the client to reconnect and if not store their data )
+ Encryption of data over socket (concurrency likely here encrypt/decrypt functions. goroutines that store broadcast the decrypted/encrypted data on the Connection Managers (Rooms) boradcast channel)