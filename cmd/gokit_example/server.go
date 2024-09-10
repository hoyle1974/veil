package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyle1974/veil/veil"
)

type ServerFactory struct {
	mux *http.ServeMux
}

func (c ServerFactory) GetServer() any {
	return c.mux
}

func server() {
	fmt.Println("-- server --")
	addr := ":8181"

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Makes sure all the server components are initialized
	// And then they will be stiched to the services being served below
	veil.VeilInitServer(ServerFactory{mux: mux})

	// Make this visible remotely
	if err := veil.Serve(&BarService{}); err != nil {
		panic(err)
	}

	// Gokit is using http transport so start up a listener
	fmt.Println("Starting server on ", addr)
	log.Fatal(server.ListenAndServe())
}
