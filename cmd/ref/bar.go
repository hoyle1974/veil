//go:generate veil
package main

import (
	"context"
	"fmt"
)

// @d:service
type Bar struct {
}

func (f *Bar) DoSomething(ctx context.Context, value int64, name string, looks []any) (string, error) {
	return fmt.Sprintf("DoSomething!:%v:%s:%v", value, name, looks), nil
}
