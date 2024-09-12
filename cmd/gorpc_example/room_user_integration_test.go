package main

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/hoyle1974/veil/veil"
	"github.com/keegancsmith/rpc"
)

type MockConnectionFactory struct{}

func (c MockConnectionFactory) GetConnection() any {
	return getRPCConn()
}

type MockServerFactory struct {
	server *rpc.Server
}

func (c MockServerFactory) GetServer() any {
	return c.server
}

func TestUserRoomIntegration(t *testing.T) {
	initGob()

	server := rpc.NewServer()

	// Start a TCP listener for the net/rpc service
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		t.Error(err)
		return
	}

	// Accept connections and serve requests
	go func() {
		for {
			fmt.Println("Waiting for a connection")
			conn, err := listener.Accept()
			if err != nil {
				// t.Error(err)
				return
			}
			fmt.Println("Connection received")
			go server.ServeConn(conn)
		}
	}()

	control := veil.InitTestFramework(MockConnectionFactory{}, MockServerFactory{server: server})
	control.StartTest(t)

	veil.Serve(NewRoomService())
	veil.Serve(NewUserService())

	// Return a function to teardown the test
	defer func(t testing.TB) {
		listener.Close()
		control.StopTest(t)
	}(t)

	roomService, err := veil.Lookup[RoomService_Interface]()
	if err != nil {
		t.Error(err)
	}
	userService, err := veil.Lookup[UserService_Interface]()
	if err != nil {
		t.Error(err)
	}

	jack, err := userService.NewUser(context.Background(), "Jack")
	if err != nil {
		t.Error(err)
	}

	jill, err := userService.NewUser(context.Background(), "Jill")
	if err != nil {
		t.Error(err)
	}

	if jack == jill {
		t.Error("User ids should be different")
	}

	roomService.AddUser(context.Background(), "Main", jill)
	roomService.AddUser(context.Background(), "Main", jack)

	users, err := roomService.GetUsers(context.Background(), "Main")
	if err != nil {
		t.Error(err)
	}
	if len(users) != 2 {
		t.Error("Not enough users")
	}

	roomService.Broadcast(context.Background(), "Main", "Hello All!")

	msg, ok, err := userService.GetLastSay(context.Background(), jack)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("no last say")
	}
	if msg != "Hello All!" {
		t.Error("Message was ", msg)
	}

	msg, ok, err = userService.GetLastSay(context.Background(), jill)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("no last say")
	}
	if msg != "Hello All!" {
		t.Error("Message was ", msg)
	}

}
