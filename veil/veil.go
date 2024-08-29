// Package veil provides functions for library users
package veil

import (
	"errors"
	"fmt"
	"reflect"
)

// Should return an opaque type that can be cast and used by a client
// to connect to the RPC service
type ConnectionFactory interface {
	GetConnection() any
}

// Clients get a client connection factory
type ClientInit func(factory ConnectionFactory)
type ServerInit func()

var clientInits = []ClientInit{}
var serverInits = []ServerInit{}

// INTERNAL - Register code to be run when client library is initialized
func RegisterClientInit(i ClientInit) {
	clientInits = append(clientInits, i)
}

// INTERNAL - Register code to be run when server library is intiialized
func RegisterServerInit(i ServerInit) {
	serverInits = append(serverInits, i)
}

// Clients should call this before client calls are executed.
func VeilInitClient(factory ConnectionFactory) {
	for _, i := range clientInits {
		i(factory)
	}
}

// Servers should call this to start the server
func VeilInitServer() {
	for _, i := range serverInits {
		i()
	}
}

var serviceRegistry = []ServiceRegistry{}

type ServiceRegistry interface {
	RPC_Bind_Service(service any) error
}

// INTERNAL - Used to register generated RPC services
func RegisterService(service ServiceRegistry) {
	serviceRegistry = append(serviceRegistry, service)
}

// Library users should call this with an instantiated version of their service
// that they want to be exposed remotely.
func Serve(service any) error {
	match := false
	for _, sr := range serviceRegistry {
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

var remoteImpl = []any{}

// INTERNAL - Registers remote implementations for client use
func RegisterRemoteImpl(service any) {
	remoteImpl = append(remoteImpl, service)
}

// Can be used to lookup a remote implementation by interface type
func Lookup[T any]() (T, error) {
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()
	for _, item := range remoteImpl {
		itemType := reflect.TypeOf(item)
		if itemType.AssignableTo(interfaceType) || itemType.Implements(interfaceType) {
			return item.(T), nil
		}
	}

	var t T
	return t, fmt.Errorf("unknown service")
}
