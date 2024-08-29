package veil

import (
	"errors"
	"fmt"
	"reflect"
)

type ClientInit func(factory ConnectionFactory)
type ServerInit func()

var clientInits = []ClientInit{}
var serverInits = []ServerInit{}

func RegisterClientInit(i ClientInit) {
	clientInits = append(clientInits, i)
}
func RegisterServerInit(i ServerInit) {
	serverInits = append(serverInits, i)
}

type ConnectionFactory interface {
	GetConnection() any
}

func VeilInitClient(factory ConnectionFactory) {
	for _, i := range clientInits {
		i(factory)
	}
}

func VeilInitServer() {
	for _, i := range serverInits {
		i()
	}
	go StartServices()
}

func StartServices() {
}

var serviceRegistry = []ServiceRegistry{}

type ServiceRegistry interface {
	RPC_Set_Service(service any) error
}

func RegisterService(service ServiceRegistry) {
	serviceRegistry = append(serviceRegistry, service)
}

func Serve(service any) error {
	match := false
	for _, sr := range serviceRegistry {
		err := sr.RPC_Set_Service(service)
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

func RegisterRemoteImpl(service any) {
	remoteImpl = append(remoteImpl, service)
}

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

func NilGet[T any](a any) T {
	var zero T
	if a == nil {
		return zero
	}
	return a.(T)
}
