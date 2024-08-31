//go:generate veil
package main

import (
	"context"
	"errors"
	"fmt"
)

// v:service
type BarService struct{}

func (f *BarService) SaySomething(ctx context.Context, name string, value int) (string, error) {
	ret := fmt.Sprintf("Hi %s, your value was %d", name, value)
	if value == -1 {
		return "", errors.New("you wanted an error")
	}
	return ret, nil
}
