package main

import (
	"context"
	"net/rpc"

	"github.com/hoyle1974/veil/veil"
)

// Registers the remote implementation with veil
// so when Lookup is called it can be returned
// It will then make a call to the remote version of the service
func VeilInitClient() {
	veil.RegisterRemoteImpl(&FooRemoteImpl{})
}

// ------------- This would be generated code to exist on the client to make a call to the server
type FooRemoteImpl struct {
}

// Make a remote call to main.Foo.Beep
func (f *FooRemoteImpl) Beep(ctx context.Context, value int) (string, error) {
	// Dial the remote RPC server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		return "", err
	}
	defer client.Close()

	request := veil.Request{
		Service: "main.Foo",
		Method:  "Beep",
		Args:    []any{value},
	}

	reply := []any{}

	err = client.Call("MyService.MyCall", request, &reply)
	if err != nil {
		return "", err
	}

	var result0 error = veil.NilGet[error](reply[0])
	var result1 string = veil.NilGet[string](reply[1])

	return result1, result0
}
