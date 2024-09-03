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

# How to write your own templates 

The provided templates can be used as examples of how to make custom templates for your own RPC needs:

* [net/rpc](https://github.com/hoyle1974/veil/blob/main/cmd/veil/rpc_service.tmpl)
* [gokit](https://github.com/hoyle1974/veil/blob/main/cmd/veil/gokit_service.tmpl)

These are written using Go's ```text/template``` package.  Veil will parse your code during the ```go generate``` command, collect all the annotated structs and then pass the data to the template function.

The data it collects will be loaded in the datastructures defined in [model](https://github.com/hoyle1974/veil/blob/main/cmd/veil/models.go)

There is still much more work needed for configurability but this is the initial prototype.

