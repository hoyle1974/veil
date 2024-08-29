# veil
Veil simplifies RPC code while trying to not be "magical".  It consist of two components.

* a code generator that uses templates to generate bindings between a Go struct and a RPC library.  By default it binds to net/RPC but I have plans to support others in the future.  
* a library used to register structs to be exposed via RPC and calls to lookup client implementations that make calls using that RPC to the exposed service.

A user might write a struct that behaves like a service like this:

```
// @d:service
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

