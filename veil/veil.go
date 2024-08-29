package veil

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"reflect"
	"sync/atomic"
)

type Init func()

var clientInits = []Init{}
var serverInits = []Init{}

func RegisterClientInit(i Init) {
	clientInits = append(clientInits, i)
}
func RegisterServerInit(i Init) {
	serverInits = append(serverInits, i)
}

func VeilInitClient() {
	for _, i := range clientInits {
		i()
	}
}

func VeilInitServer() {
	for _, i := range serverInits {
		i()
	}
	go StartServices()
}

var conn atomic.Pointer[rpc.Client]

func newConn() (*rpc.Client, error) {
	return rpc.Dial("tcp", "localhost:1234")
}

func GetConn() *rpc.Client {
	if conn.Load() != nil {
		return conn.Load()
	}

	db, err := newConn()
	if err != nil {
		panic(err)
	}

	old := conn.Swap(db)
	if old != nil {
		old.Close()
	}
	return db
}

func StartServices() {

	// Start a TCP listener
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}

	// Accept connections and serve requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

var serviceRegistry = []ServiceRegistry{}

type ServiceRegistry interface {
	RPC_Set_Service(service any) error
}

func RegisterService(service ServiceRegistry) {
	fmt.Println("RegisterService: ", reflect.TypeOf(service))

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
