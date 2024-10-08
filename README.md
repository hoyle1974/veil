# veil
Veil simplifies RPC code while trying to not be "magical".  It consist of two components.

* a code generator that uses templates to generate bindings between a Go struct and an RPC library.  By default it binds to net/RPC but I have plans to support others in the future.  
* a library used to register structs to be exposed via RPC and calls to lookup client implementations that make calls using that RPC to the exposed service.

A user might write a struct that behaves like a service like this:

```
// @v:service -t rpc
type Foo struct {
}

func (r *RoomService) MyCall(ctx context.Context, name string, value int) (string, error) {
    ret := fmt.Sprintf("Hello %s, your value is %d",name,value)
    if ret==-1 {
        return "", errors.New("Here is your error!")
    }
    return ret, nil
}
```

When you run ```go generate``` on your code it will find all structs with a @v:service comment and work on exposing all properly defined methods on that struct.  The method call must be exported (upper case name), the first argument must be a context.Context and the last return value will be an error.

Currently veil supports net/rpc and gokit, these can be configured using the -t flag which can either be rpc (for net/rpc), gokit, or a path to a template of your choosing.  Instead of defining this in your code you can put these options in a config file (in either ~/.veil or at the location defined in the environment variable VEIL_CONFIG_FILE).  You can also set the environment variable VEIL_CONFIG directly.

Using Veil they would expose this struct via RPC like this:

```
veil.VeilInitServer()

// Make Foo visible remotely
if err := veil.Serve(&Foo{}); err != nil {
	panic(err)
}

// Setup net/rpc like normal
// Start a TCP listener
listener, err := net.Listen("tcp", ":1234")
if err != nil {
	fmt.Println("Listen error:", err)
	return
}

// Accept connections and serve requests
for {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("error:", err)
		continue
	}
	go rpc.ServeConn(conn)
}
```

Now from a client they could do this

```
// Connection factory is a user supplied interface GetConnection() method that returns an opaque
// connection for client calls.  In this case it would return a *rpc.Client
veil.VeilInitClient(ConnFactory{})

// Create a context, because this will be a network call
ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second))

// Lookup a client to talk to our service with the generated interface
foo, err := veil.Lookup[Foo_Interface]()
if err != nil {
	panic(err)
}

// Now we can call the remote function
ret, err := foo.MyCall(ctx, "Jack", 5)
if err != nil {
    panic(err)
}
fmt.Println(ret)

// This call will return the error
_, err := foo.MyCall(ctx, "Jack", -1)
fmt.Println(err)

```

All the code to bind your exposed struct to the chosen RPC is generated automatically.  The code necessary to call your RPC is also generated.  Network errors are automatically stitched into the returned error value as needed.  Context deadlines & cancels should work across the network as well.

Right now nothing else is added as part of the generated code.  My plan is to play with this in a personal project for a bit, see how it behaves and performs, and make adjustments as needed.  Feedback is appreciated.

**Next Steps**

This is simply a prototype to explore the design.  Next steps will include:

* More configurability
* unit tests
    * In the library itself
    * generated helper code to test client & server locally with mock networks
* Add support for:
    * other RPC libraries like grpc
    * authentication, authorization
    * service discovery 
    * rate limiting
    * circuit breaker 