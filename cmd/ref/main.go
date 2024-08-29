package main

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync/atomic"
	"time"

	"github.com/hoyle1974/veil/veil"
)

var conn atomic.Pointer[rpc.Client]

func newConn() (*rpc.Client, error) {
	return rpc.Dial("tcp", "localhost:1234")
}

func GetRPCConn() *rpc.Client {
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

func server() {
	fmt.Println("server.")

	veil.VeilInitServer()

	// Make these visible remotely
	if err := veil.Serve(&RoomService{}); err != nil {
		panic(err)
	}
	if err := veil.Serve(&UserService{}); err != nil {
		panic(err)
	}

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

func client() {
	veil.VeilInitClient()

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second))

	// Lookup the remote interface for the user service
	// then create 2 users
	users, err := veil.Lookup[UserService_Interface]()
	if err != nil {
		panic(err)
	}
	userId1, err := users.NewUser(ctx, "Joe")
	if err != nil {
		panic(err)
	}
	userId2, err := users.NewUser(ctx, "Bob")
	if err != nil {
		panic(err)
	}

	// Lookup the remote interface for the room service
	// Add the users
	// Then broadcast a message to all users
	roomId := "General Lobby"
	rooms, err := veil.Lookup[RoomService_Interface]()
	if err != nil {
		panic(err)
	}
	if _, _, err = rooms.AddUser(ctx, roomId, userId1); err != nil {
		panic(err)
	}
	if _, _, err = rooms.AddUser(ctx, roomId, userId2); err != nil {
		panic(err)
	}
	if _, err = rooms.Broadcast(ctx, roomId, "Hello Everyone!"); err != nil {
		panic(err)
	}
}

// func main() {

// }

func main() {
	if os.Args[1] == "server" {
		server()
	}
	if os.Args[1] == "client" {
		client()
	}

	fmt.Println("Waiting.")
	select {}
}
