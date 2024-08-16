package main

import (
	"context"
	"fmt"
)

//go:generate veil

// @d:service
type Foo struct {
}

func (f *Foo) Beep(ctx context.Context) error {
	return nil
}

func (f *Foo) BeepBad() error {
	return nil
}

// @d:service
type Bar struct {
}

func (f *Bar) Boop(ctx context.Context) {

}

func (f *Bar) BoopBad(ctx context.Context) string {
	return ""
}

func main() {
	a := &Foo{}
	var b FooInterface

	b = a

	fmt.Println(a)
	fmt.Println(b)
}
