package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hoyle1974/veil/veil"
)

//go:generate veil

// @d:service
type Foo struct {
}

func (f *Foo) Beep(ctx context.Context, value int) (string, error) {
	fmt.Println("-- beep --")
	return fmt.Sprintf("beep!:%v", value), nil
}

func startServer() {
	go veil.Serve(&Foo{})
}

func startClient() {
	// This is the client
	foo, err := veil.Lookup[FooInterface]()
	if err != nil {
		panic(err)
	}

	s, err := foo.Beep(context.Background(), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Beep Returned", s)
}

func main() {

	VeilInitClient()
	VeilInitServer()

	startServer()

	time.Sleep(time.Second)

	startClient()

}
