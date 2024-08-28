package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hoyle1974/veil/veil"
)

func server() {
	fmt.Println("server.")

	veil.VeilInitServer()

	// Make these visible remotely
	veil.Serve(&RoomService{})
	veil.Serve(&UserService{})

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
