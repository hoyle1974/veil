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

	args := []any{}
	args = append(args, value)

	request := veil.Request{
		Service: "main.Foo",
		Method:  "Beep",
		Args:    args,
	}

	reply := []any{}

	err = client.Call("MyService.MyCall", request, &reply)
	if err != nil {
		return "", err
	}

	result0 := reply[0].(string)

	return result0, nil
}
