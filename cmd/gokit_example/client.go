package main

import (
	"context"
	"fmt"

	"github.com/hoyle1974/veil/veil"
)

// The client will need a connection factory, in this case
// We just provide the url for the connection
type ConnFactory struct{}

func (c ConnFactory) GetConnection() any {
	return "http://localhost:8181"
}

func client() {
	fmt.Println("-- client --")

	// Makes sure all the client components are initialized
	// And they have access to the connection factory
	veil.VeilInitClient(ConnFactory{})

	// Lookup the interface we want
	bar, err := veil.Lookup[BarService_Interface]()
	if err != nil {
		panic(err)
	}

	// Call the method, this is an RPC over HTTP under the hood
	result, err := bar.SaySomething(context.Background(), "Jack", 2431)
	if err != nil {
		panic(err)
	}

	// Show the results
	fmt.Println("Result:", result)
}
