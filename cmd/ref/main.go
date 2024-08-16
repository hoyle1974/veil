package main

import (
	"context"

	"github.com/hoyle1974/veil/veil"
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
