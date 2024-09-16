package main

import (
	"context"
	"fmt"

	"github.com/hoyle1974/veil/veil"
)

func client() {
	fmt.Println("-- client --")

	// Makes sure all the client components are initialized
	// And they have access to the connection factory
	veil.VeilInitClient(veil.GetLocalConnectionFactory())

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
