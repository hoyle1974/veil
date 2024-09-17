package main

import "context"

// Start the server first: go run . server
// Then run the client: go run . client
func main() {

	go server()

	client()

	b := &BarService{}
	b.Ping.PongPing(context.Background())
	// b.Ping.Ping2.PingPong(context.Background())
	// b.Ping.PingPong(context.Background())
	b.PingPong(context.Background())

}
