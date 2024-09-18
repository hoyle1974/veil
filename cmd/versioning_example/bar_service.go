//go:generate veil
package main

import (
	"context"
	"errors"
	"fmt"
)

// v:service -t local
type BarService struct {
	Version string
}

// v:service -t local
type BarService2 struct {
	BarService
	Version string
}

func (f *BarService) SaySomething(ctx context.Context, name string, value int) (string, error) {
	ret := fmt.Sprintf("Hi %s, your value was %d, my version is %s", name, value, f.Version)
	if value == -1 {
		return "", errors.New("you wanted an error")
	}
	return ret, nil
}

func (f *BarService2) SaySomething(ctx context.Context, name string, value int, extra string) (string, error) {
	ret := fmt.Sprintf("Hi %s, your value was %d, my extra is %s, my version is %s", name, value, extra, f.Version)
	if value == -1 {
		return "", errors.New("you wanted an error")
	}
	return ret, nil
}
