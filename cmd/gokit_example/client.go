package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/hoyle1974/veil/veil"
)

// The client will need a connection factory, in this case
// We just provide the url for the connection
type ConnFactory struct{}

func (c ConnFactory) GetConnection() any {
	return Connection{}
}

type Connection struct{}

func (c Connection) Get(path string, jsonData []byte) (*http.Response, error) {
	url := "http://localhost:8181" + path
	fmt.Println("Get ", url)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// Make the HTTP request
	return client.Do(req)
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
