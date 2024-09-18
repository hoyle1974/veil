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
	service := &BarService2{BarService: BarService{Version: "1"}, Version: "2"}
	if err := veil.Serve(&service.BarService); err != nil {
		panic(err)
	}
	if err := veil.Serve(service); err != nil {
		panic(err)
	}
}
