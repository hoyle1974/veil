package main

import "fmt"

//go:generate veil

// @d:service
type Foo struct {
}

func (f *Foo) Beep() error {
	return nil
}

// @d:service
type Bar struct {
}

func (f *Bar) Boop() {

}

func main() {
	a := &Foo{}
	var b FooInterface

	b = a

	fmt.Println(a)
	fmt.Println(b)
}
