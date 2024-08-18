package veil

import (
	"fmt"
	"net"
	"net/rpc"
	"reflect"
	"sync"
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

func getConn() *rpc.Client {
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

func Call(request Request, reply *[]any) error {
	conn, err := newConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Call("MyService.MyCall", request, &reply)
	if err != nil {
		return err
	}
	return nil
}

type Request struct {
	Service string
	Method  string
	Args    []byte
}

type MyService struct {
}

type RPCCB func(any, string, any, *[]any)

var slock sync.Mutex
var services = map[string]any{}
var cbs = map[string]RPCCB{}

func RegisterService(service string, cb RPCCB) {
	cbs[service] = cb
}

func (t *MyService) MyCall(request *Request, reply *[]any) error {
	cbs[request.Service](services[request.Service], request.Method, request.Args, reply)
	return nil
}

func StartServices() {
	// Register the Arithmetic type with the RPC server
	rpc.Register(&MyService{})

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

func Serve(service any) {
	slock.Lock()
	defer slock.Unlock()

	name := reflect.TypeOf(service).String()[1:]
	services[name] = service
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
