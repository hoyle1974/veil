package main

// Registers the remote implementation with veil
// so when Lookup is called it can be returned
// It will then make a call to the remote version of the service
// func VeilInitClient() {
// 	veil.RegisterRemoteImpl(&FooRemoteImpl{})
// }

/*
// ------------- This would be generated code to exist on the client to make a call to the server
type FooRemoteImpl struct {
}

// Make a remote call to main.Foo.Beep
func (f *FooRemoteImpl) Beep(ctx context.Context, value int) (string, error) {

	request := veil.Request{
		Service: "main.Foo",
		Method:  "Beep",
		Args:    []any{value},
	}

	reply := []any{}
	var result0 error
	var result1 string

	err := veil.Call(request, &reply)
	if err != nil {
		result0 = err
	} else {
		result0 = veil.NilGet[error](reply[0])
		result1 = veil.NilGet[string](reply[1])
	}

	return result1, result0
}

// Make a remote call to main.Foo.Beep
func (f *FooRemoteImpl) Boop(ctx context.Context, value string) (string, error) {

	request := veil.Request{
		Service: "main.Foo",
		Method:  "Boop",
		Args:    []any{value},
	}

	reply := []any{}
	var result0 error
	var result1 string

	err := veil.Call(request, &reply)
	if err != nil {
		result0 = err
	} else {
		result0 = veil.NilGet[error](reply[0])
		result1 = veil.NilGet[string](reply[1])
	}

	return result1, result0
}
*/
