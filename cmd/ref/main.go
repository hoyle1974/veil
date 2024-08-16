package main

import (
	"context"
)

//go:generate veil

// @d:service
type Foo struct {
}

func (f *Foo) Beep(ctx context.Context) (string, error) {
	return "beep!", nil
}

func main() {
	veil.Serve(&Foo{})
}
