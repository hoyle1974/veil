package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func initGob() {
	gob.Register(UserId(""))
	gob.Register(RoomId(""))
	gob.Register([]UserId{})
	gob.Register([]RoomId{})
}

// Start the server first: go run . server
// Then run the client: go run . client
func main() {
	initGob()

	if len(os.Args) == 1 || os.Args[1] == "server" {
		server()

		fmt.Println("Waiting.")
		select {}
	}
	if os.Args[1] == "client" {
		client()
	}

}
