### Design


### Setting errors on response writier
+ https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter


## Tests
### Todo
1. Write tests to veryify the clients are being created correctly in the correct rooms(hubs/managers/broker)
    1. __ACTIVE__ Test1: Test that two requests to the same chatid creates only one new entry in our MAP
        + Removed race condition BUT still have to figure out how to upgrade the connection to a socket 
        ``` panic: websocket: response does not implement http.Hijacker ```
        + Make own websocket client also disect the code from the first bullet and figure out what it is doing (in relation to testing a websocket connection)
            + https://github.com/posener/wstest/blob/master/dialer.go
            + https://godoc.org/github.com/gorilla/websocket#Dialer
            + https://stackoverflow.com/questions/32745716/i-need-to-connect-to-an-existing-websocket-server-using-go-lang
    + Test2: Check that two concurrent requests to the same chatID creates only one new entry in our MAP
    + Test3: Check that two concurrent exits from the same chat Does not cause a race condition (writing to brokermap)

- Stop the broker if all clients have left and remove it from the map so it doesn't stay running for no reason (needs testing)

## General Todo
+ __Implement go client for connecting to the chat server__ https://stackoverflow.com/questions/32745716/i-need-to-connect-to-an-existing-websocket-server-using-go-lang
+ Add os signal package to our racerd and use context with each route, this way we can gracefully shutdown the program if something happens 
    + Put server on a select with a signal channel
    + If the server fails cancel parent context, if the signal is canceled cancel parent context 


### Links
1. Testing chi router https://github.com/go-chi/chi/issues/76
    1. Accessing map concurrently in go https://stackoverflow.com/questions/52512915/how-to-solve-concurrency-access-of-golang-map
    2. Read Write Mutex in depth https://stackoverflow.com/questions/19148809/how-to-use-rwmutex-in-golang
+ Buffers & IO Pipes explained - https://medium.com/stupid-gopher-tricks/streaming-data-in-go-without-buffering-3285ddd2a1e5


+ Work on implementing client interface to decouple brokers from clients 
+ Figure out what time variables are doing in client.go




+ transmit json data through the socket
+ Make Room initialization happen inside a handler function
+ Add database support for retreival of chats for a given id 
+ research concurrent design patterns 
+ research backing up data on client interruption (how can I block for a certain period of time and wait for the client to reconnect and if not store their data )
+ Encryption of data over socket (concurrency likely here encrypt/decrypt functions. goroutines that store broadcast the decrypted/encrypted data on the Connection Managers (Rooms) boradcast channel)