package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hoyle1974/veil/veil"
)

func startServer() {
	go veil.Serve(&Foo{})
	go veil.Serve(&Bar{})
}

func startClient() {
	// This is the client
	foo, err := veil.Lookup[FooInterface]()
	if err != nil {
		panic(err)
	}

	// bar, err := veil.Lookup[BarInterface]()
	// if err != nil {
	// 	panic(err)
	// }

	// s, err := foo.Beep(context.Background(), 5)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Beep Returned", s)

	s, err := foo.Boop(context.Background(), "Hello World")
	if err != nil {
		panic(err)
	}
	fmt.Println("Boop Returned", s)

	// s, err = bar.DoSomething(context.Background(), 1, "BarBar", []any{1, 2, 3, 4})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Bar Returned", s)
}

func main() {
	veil.VeilInitClient()
	veil.VeilInitServer()

	startServer()

	time.Sleep(time.Second)

	startClient()

}
