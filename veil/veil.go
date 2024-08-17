package veil

import (
	"fmt"
	"net"
	"net/rpc"
	"reflect"
)

type Request struct {
	Service string
	Method  string
	Args    []any
}

type MyService struct {
	name    string
	service any
}

type RPCCB func(any, string, []any, *[]any)

var services = map[string]*MyService{}
var cbs = map[string]RPCCB{}

func RegisterService(service string, cb RPCCB) {
	cbs[service] = cb
}

func (t *MyService) MyCall(request *Request, reply *[]any) error {
	cbs[request.Service](t.service, request.Method, request.Args, reply)
	return nil
}

func Serve(service any) {
	name := reflect.TypeOf(service).String()[1:]
	s := &MyService{name: name, service: service}
	services[name] = s

	// Register the Arithmetic type with the RPC server
	rpc.Register(s)

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

var l = []any{}

func RegisterRemoteImpl(service any) {
	l = append(l, service)
}

func Lookup[T any]() (T, error) {
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()
	for _, item := range l {
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
