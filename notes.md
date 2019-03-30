### Design
+ Add database support for retreival of chats for a given id 
+ research backing up data on client interruption (how can I block for a certain period of time and wait for the client to reconnect and if not store their data )
+ Encryption of data over socket (concurrency likely here encrypt/decrypt functions.NewBrokers (Rooms) boradcast channel)
+ Seperate rooms types for group and private chats
    + Two types of managers private or public, in public any client can join.
    + in private only clients that have been whitelisted can join



### Formating time
+ Good info about how to translate to different local timezones in go https://stackoverflow.com/questions/25318154/convert-utc-to-local-time-go
+ https://unix4lyfe.org/time/?v=1
+ https://stackoverflow.com/questions/9539108/obtaining-a-unix-timestamp-in-go-language-current-time-in-seconds-since-epoch


## Handling http requests with handlers
### Testing
+ example hashicorp
    + https://github.com/hashicorp/vault/blob/master/http/handler_test.go
    + https://github.com/hashicorp/vault/blob/master/http/testing.go
    + https://github.com/hashicorp/go-cleanhttp
+ example benjonshon
    + https://github.com/benbjohnson/pilosa/blob/master/handler.go
    + https://github.com/benbjohnson/pilosa/blob/master/handler_test.go

## export_test and testing unexported functions from packages 
+ https://www.reddit.com/r/golang/comments/7x1dw5/some_thoughts_on_testing_in_go/
+ https://www.reddit.com/r/golang/comments/a2ysk8/tests_in_test_package_or_not/



## Crypto & rand & unique identifiers in go
+ http://blog.sgmansfield.com/2016/01/the-hidden-dangers-of-default-rand/
+ https://stackoverflow.com/questions/2662591/go-example-about-crypto-rand



## Select statements and two or more ready channels
+ https://stackoverflow.com/questions/47645808/how-does-select-work-when-multiple-channels-are-involved


## Noop read closer and reading from request body
+ https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361


## Good info on how to use the context package correctly
+ https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/
+ https://www.reddit.com/r/golang/comments/6s5a7q/context_should_go_away_for_go_2/?utm_source=reddit-android


## Errors & solutions
### GCC compiler error using go modules and go test -race ./...
+ https://github.com/golang/go/issues/27089 -- offers a solution
+ https://github.com/golang/go/issues/26988
### Difference between compilers 
+ https://stackoverflow.com/questions/25811445/what-are-the-primary-differences-between-gc-and-gccgo


### Htttp & Server
+ There are a number of issues to consider when using golang as an http server
1. __Headers__
    + https://www.reddit.com/r/golang/comments/7yctil/which_http_headers_should_i_include_in_my_api/
2. __Read & Write Timeouts__ (not setting these can cause leaky goroutines keeping connections alive way longer than necessary)
    + https://stackoverflow.com/questions/10971800/golang-http-server-leaving-open-goroutines/10972453#10972453
    + https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779 (read article and comments good info)


### Sockets
+ ``` panic: websocket: response does not implement http.Hijacker ```
+ https://github.com/posener/wstest/blob/master/dialer.go
+ https://godoc.org/github.com/gorilla/websocket#Dialer
+ https://stackoverflow.com/questions/32745716/i-need-to-connect-to-an-existing-websocket-server-using-go-lang
+ Cannont use the same endpoint to connect via sock and http https://stackoverflow.com/questions/48006498/is-this-possible-to-server-websocket-handler-and-normal-servlet-over-same-contex



### Pointers and values
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




## General Links
1. Testing chi router https://github.com/go-chi/chi/issues/76
    1. Accessing map concurrently in go https://stackoverflow.com/questions/52512915/how-to-solve-concurrency-access-of-golang-map
    2. Read Write Mutex in depth https://stackoverflow.com/questions/19148809/how-to-use-rwmutex-in-golang
+ Buffers & IO Pipes explained - https://medium.com/stupid-gopher-tricks/streaming-data-in-go-without-buffering-3285ddd2a1e5
    + bytes buffers and when to use them and how they alloc mem https://syslog.ravelin.com/bytes-buffer-i-thought-you-were-my-friend-4148fd001229



### Setting errors on response writier
+ https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter