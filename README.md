# veil
Veil simplifies RPC code while trying to not be "magical".  It consist of two components.

* a code generator that uses templates to generate bindings between a Go struct and an RPC library.  By default it binds to net/RPC but I have plans to support others in the future.  
* a library used to register structs to be exposed via RPC and calls to lookup client implementations that make calls using that RPC to the exposed service.

A user might write a struct that behaves like a service like this:

```
// @v:service
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

When you run ```go generate``` on your code it will find all structs with a @v:service comment and work on exposing all properly defined methods on that struct.  The method call must be exported (upper case name), the first argument must be a context.Context and the last reurn value will be an error.

Using Veil they would expose this struct via RPC like this

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

**Next Steps**

This is simply a prototype to explore the design.  Next steps will include:

* Support for other RPC paradigms, maybe grpc
* More configurability
** generated struct names 
** generated file names