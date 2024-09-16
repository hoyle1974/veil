// Package veil provides functions for library users
package veil

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hoyle1974/veil/veil_internal"
)

type localConnection struct {
}

func (l localConnection) Get(name string) any {
	return services[name]
}

type localConnectionFactory struct {
}

func (l localConnectionFactory) GetConnection() any {
	return localConnection{}
}

func GetLocalConnectionFactory() ConnectionFactory {
	return localConnectionFactory{}
}

// Should return an opaque type that can be cast and used by a client
// to connect to the RPC service
type ConnectionFactory interface {
	GetConnection() any
}

type ServerFactory interface {
	GetServer() any
}

// Clients get a client connection factory
type ClientInit func(factory ConnectionFactory)
type ServerInit func(factory ServerFactory)

// Clients should call this before client calls are executed.
func VeilInitClient(factory ConnectionFactory) {
	for _, i := range veil_internal.ClientInits {
		i(factory)
	}
}

// Servers should call this to start the server
func VeilInitServer(factory ServerFactory) {
	for _, i := range veil_internal.ServerInits {
		i(factory)
	}
}

// Library users should call this with an instantiated version of their service
// that they want to be exposed remotely.
var services = map[string]any{}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func Serve(service any) error {
	match := false
	for _, sr := range veil_internal.ServiceRegistry {
		err := sr.RPC_Bind_Service(service)
		if err == nil {
			match = true
			services[getType(service)] = service
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
