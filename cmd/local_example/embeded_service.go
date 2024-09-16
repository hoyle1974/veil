package main

import (
	"context"
	"errors"
	"fmt"
)

func (f *BarService) SaySomethingElse(ctx context.Context, name string, value int) (string, error) {
	ret := fmt.Sprintf("Hi %s, your value was %d - ELSE", name, value)
	if value == -1 {
		return "", errors.New("you wanted an error")
	}
	return ret, nil
}

type Ping struct {
}

func (f *Ping) Ping(ctx context.Context) (string, error) {
	return "Pong", nil
}
