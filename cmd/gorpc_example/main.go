package main

import (
	"fmt"
	"os"
)

// Start the server first: go run . server
// Then run the client: go run . client
func main() {

	if len(os.Args) == 1 || os.Args[1] == "server" {
		server()

		fmt.Println("Waiting.")
		select {}
	}
	if os.Args[1] == "client" {
		client()
	}

}
