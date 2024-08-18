//go:generate veil
package main

import (
	"context"
	"fmt"
)

// @d:service
type Foo struct {
}

func (f *Foo) Beep(ctx context.Context, value int) (string, error) {
	return fmt.Sprintf("beep!:%v", value), nil
}

func (f *Foo) Boop(ctx context.Context, value string) (string, error) {
	return fmt.Sprintf("boop!:%v", value), nil
}
