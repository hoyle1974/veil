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
	bar1, err := veil.Lookup[BarService_Interface]()
	if err != nil {
		panic(err)
	}

	// Lookup the interface we want
	bar2, err := veil.Lookup[BarService2_Interface]()
	if err != nil {
		panic(err)
	}

	// Call the method, this is an RPC over HTTP under the hood
	result, err := bar1.SaySomething(context.Background(), "Jack", 123)
	if err != nil {
		panic(err)
	}
	// Show the results
	fmt.Println("Result:", result)

	// Call the method, this is an RPC over HTTP under the hood
	result, err = bar2.SaySomething(context.Background(), "Jill", 456, "my extra!")
	if err != nil {
		panic(err)
	}

	// Show the results
	fmt.Println("Result:", result)
}
