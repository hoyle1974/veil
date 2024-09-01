package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "server" {
		server()
	}
	if os.Args[1] == "client" {
		client()
	}

	fmt.Println("Waiting.")
	select {}
}
