### Design


+ Add database support for retreival of chats for a given id 
+ research backing up data on client interruption (how can I block for a certain period of time and wait for the client to reconnect and if not store their data )
+ Encryption of data over socket (concurrency likely here encrypt/decrypt functions. goroutines that store broadcast the decrypted/encrypted data on the Connection Managers (Rooms) boradcast channel)

## Tests
### Todo
1. Done
 + Write tests to veryify the clients are being created correctly in the correct rooms(hubs/managers/broker)
    1. _Test1: Test that two requests to the same chatid creates only one new entry in our MAP
        + Removed race condition BUT still have to figure out how to upgrade the connection to a socket 
        ``` panic: websocket: response does not implement http.Hijacker ```
        + Make own websocket client also disect the code from the first bullet and figure out what it is doing (in relation to testing a websocket connection)
            + https://github.com/posener/wstest/blob/master/dialer.go
            + https://godoc.org/github.com/gorilla/websocket#Dialer
            + https://stackoverflow.com/questions/32745716/i-need-to-connect-to-an-existing-websocket-server-using-go-lang
            + Cannont use the same endpoint to connect via sock and http https://stackoverflow.com/questions/48006498/is-this-possible-to-server-websocket-handler-and-normal-servlet-over-same-contex
    + Test2: Check that two concurrent requests to the same chatID creates only one new entry in our MAP
    + Test3: Check that two concurrent exits from the same chat Does not cause a race condition (writing to brokermap)


## General Todo
+ Possibly Implement go client for connecting to the chat server https://stackoverflow.com/questions/32745716/i-need-to-connect-to-an-existing-websocket-server-using-go-lang
+ Add os signal package to our racerd and use context with each route, this way we can gracefully shutdown the program if something happens 
    + Put server on a select with a signal channel
    + If the server fails cancel parent context, if the signal is canceled cancel parent context 


### Links
1. Testing chi router https://github.com/go-chi/chi/issues/76
    1. Accessing map concurrently in go https://stackoverflow.com/questions/52512915/how-to-solve-concurrency-access-of-golang-map
    2. Read Write Mutex in depth https://stackoverflow.com/questions/19148809/how-to-use-rwmutex-in-golang
+ Buffers & IO Pipes explained - https://medium.com/stupid-gopher-tricks/streaming-data-in-go-without-buffering-3285ddd2a1e5
    + bytes buffers and when to use them and how they alloc mem https://syslog.ravelin.com/bytes-buffer-i-thought-you-were-my-friend-4148fd001229


## Pointers and values
+ Encountered a problem unmarshalling json and wanted to explain the solution

```
var chatmsg *message
fmt.Printf("Chat message %p\n", &chatmsg)
```

> Q: What is the value of chatmsg ?

A: The value of chatmsg is nil. Remember that everything in go is pass by value. If you pass chatmsg into json unmarshall you are passing in nil.
Unmarshall needs an underlying struct to place its data, It needs an address of a location in memory to store the values it creates.
We declared a new empty pointer to a message. The value of an empty pointer is nil. If we pass this value around it will be nil. IF we wanted this value to be populated by unmarshal, __we would have to pass ITS address__ By Passing the address of the empty pointer, go knows we want to fill that memory address with the concrete value of the pointer

It could also help to think of *Type & Type as two seperate values in go. Type is a concrete value and declaring it, it will be initialized to the 0 value for its type. Declaring *Type will intialize only a pointer (1 word) whose value is nil.


### Setting errors on response writier
+ https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter