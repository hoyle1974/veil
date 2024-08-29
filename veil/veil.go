// Package veil provides functions for library users
package veil

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hoyle1974/veil/veil_internal"
)

// Should return an opaque type that can be cast and used by a client
// to connect to the RPC service
type ConnectionFactory interface {
	GetConnection() any
}

// Clients get a client connection factory
type ClientInit func(factory ConnectionFactory)
type ServerInit func()

// Clients should call this before client calls are executed.
func VeilInitClient(factory ConnectionFactory) {
	for _, i := range veil_internal.ClientInits {
		i(factory)
	}
}

// Servers should call this to start the server
func VeilInitServer() {
	for _, i := range veil_internal.ServerInits {
		i()
	}
}

// Library users should call this with an instantiated version of their service
// that they want to be exposed remotely.
func Serve(service any) error {
	match := false
	for _, sr := range veil_internal.ServiceRegistry {
		err := sr.RPC_Bind_Service(service)
		if err == nil {
			match = true
		}
	}
	if !match {
		return errors.New("service not supported")
	}
	return nil
}

// Can be used to lookup a remote implementation by interface type
func Lookup[T any]() (T, error) {
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()
	for _, item := range veil_internal.RemoteImpl {
		itemType := reflect.TypeOf(item)
		if itemType.AssignableTo(interfaceType) || itemType.Implements(interfaceType) {
			return item.(T), nil
		}
	}

	var t T
	return t, fmt.Errorf("unknown service")
}
