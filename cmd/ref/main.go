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
	return fmt.Sprintf("beep!:%v", value), nil
}

func (f *Foo) Boop(ctx context.Context, value string) (string, error) {
	return fmt.Sprintf("boop!:%v", value), nil
}

// @d:service
type Bar struct {
}

func (f *Bar) DoSomething(ctx context.Context, value int64, name string, looks []any) (string, error) {
	return fmt.Sprintf("DoSomething!:%v:%s:%v", value, name, looks), nil
}

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

	bar, err := veil.Lookup[BarInterface]()
	if err != nil {
		panic(err)
	}

	s, err := foo.Beep(context.Background(), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Beep Returned", s)

	s, err = foo.Boop(context.Background(), "Hello World")
	if err != nil {
		panic(err)
	}
	fmt.Println("Boop Returned", s)

	s, err = bar.DoSomething(context.Background(), 1, "BarBar", []any{1, 2, 3, 4})
	if err != nil {
		panic(err)
	}
	fmt.Println("Bar Returned", s)
}

func main() {

	// data, err := json.Marshal(Bar_DoSomething_Request{5, "Name", []any{1, 2, 3, 4, 5}})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(data)

	// r := Bar_DoSomething_Request{}
	// err = json.Unmarshal(data, &r)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(r)

	VeilInitClient()
	VeilInitServer()

	startServer()

	time.Sleep(time.Second)

	startClient()

}
