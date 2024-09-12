package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hoyle1974/veil/veil"
	"github.com/keegancsmith/rpc"
)

var conn atomic.Pointer[rpc.Client]

func newConn() (*rpc.Client, error) {
	return rpc.Dial("tcp", "localhost:1234")
}

func getRPCConn() *rpc.Client {
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

type ConnFactory struct{}

func (c ConnFactory) GetConnection() any {
	return getRPCConn()
}

func client() {
	fmt.Println("-- client --")

	// Makes sure all the client components are initialized
	// And they have access to the connection factory
	veil.VeilInitClient(ConnFactory{})

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
	roomId := RoomId("General Lobby")
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
	if value, err := rooms.Broadcast(ctx, roomId, "Hello Everyone!"); err != nil {
		fmt.Println("Broadcast had a problem!")
		panic(err)
	} else {
		fmt.Println("Broadcast Value = ", value)
	}
}
