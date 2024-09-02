package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyle1974/veil/veil"
)

func server() {
	fmt.Println("-- server --")

	// Makes sure all the server components are initialized
	// And then they will be stiched to the services being served below
	veil.VeilInitServer()

	// Make this visible remotely
	if err := veil.Serve(&BarService{}); err != nil {
		panic(err)
	}

	// Gokit is using http transport so start up a listener
	log.Fatal(http.ListenAndServe(":8181", nil))

}
