package main

import (
	"fmt"
	"net"

	"github.com/hoyle1974/veil/veil"
	"github.com/keegancsmith/rpc"
)

func server() {
	fmt.Println("-- server --")

	// Makes sure all the server components are initialized
	// And then they will be stiched to the services being served below
	veil.VeilInitServer()

	// Make these visible remotely
	if err := veil.Serve(&RoomService{}); err != nil {
		panic(err)
	}
	if err := veil.Serve(&UserService{}); err != nil {
		panic(err)
	}

	// Start a TCP listener for the net/rpc service
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}

	// Accept connections and serve requests
	for {
		fmt.Println("Waiting for a connection")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		fmt.Println("Connection received")
		go rpc.ServeConn(conn)
	}
}
