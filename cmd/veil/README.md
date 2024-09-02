This is the code for the veil generator.  To use in your go files you need to install this package using ```go install``` and then add this directive at the top of the go file syou want it to parse ```//go:generate veil```.  You can then annotate your structs that you want veil to parse and write RPC binding code for.  The annotation might look like this:

```
// v:service -t gokit
type BarService struct{}

// This function matches the signature that veil will expose via RPC.
func (f *BarService) Foo(ctx context.Context, name string, value int) (string, error) {
    . . .
}
```

When complete you will find some impl_*.go files that contain:

* An interface for your struct defining the methods it will expose
* Bindings that will expose your service over the RPC of your choice.
* A struct that implements the interface and provides client methods to call your new RPC service
* Initialization code that let's the Veil system do any needed work.

You are repsonsible for the RPC server's initialization.  You can use Veil to lookup on the client side the remote implementation of your service.

Please see the examples in cmd/gokit_example and cmd/gorpc_example for complete simple examples.

TODO - Documentation on how to use Veil to generate bindings that match your own needs if the default ones do not.