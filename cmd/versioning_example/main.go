package main

// Start the server first: go run . server
// Then run the client: go run . client
func main() {

	go server()

	client()

}
