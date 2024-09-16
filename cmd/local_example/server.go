package main

import (
	"fmt"

	"github.com/hoyle1974/veil/veil"
)

type ServerFactory struct {
}

func (c ServerFactory) GetServer() any {
	return true
}

func server() {
	fmt.Println("-- server --")

	// Makes sure all the server components are initialized
	// And then they will be stiched to the services being served below
	veil.VeilInitServer(ServerFactory{})

	// Make this visible remotely
	if err := veil.Serve(&BarService{}); err != nil {
		panic(err)
	}
}
