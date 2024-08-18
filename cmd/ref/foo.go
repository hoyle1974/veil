//go:generate veil
package main

import (
	"context"
	"fmt"

	"github.com/hoyle1974/veil/veil"
)

// @d:service
type Foo struct {
}

func (f *Foo) Beep(ctx context.Context, value int) (string, error) {
	return fmt.Sprintf("beep!:%v", value), nil
}

func (f *Foo) Boop(ctx context.Context, value string) (string, error) {
	bar, err := veil.Lookup[BarInterface]()
	if err != nil {
		panic(err)
	}
	s, err := bar.DoSomething(context.Background(), 1, "BarBar", []any{1, 2, 3, 4})
	if err != nil {
		panic(err)
	}
	fmt.Println("Bar Returned", s)

	return fmt.Sprintf("boop!:%v", value), nil
}
