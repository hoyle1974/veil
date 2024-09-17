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

type Ping2 struct {
}

func (f *Ping2) PingPong(ctx context.Context) (string, error) {
	return "PongPong", nil
}

type Ping struct {
	Ping2
}

func (f *Ping) PongPing(ctx context.Context) (string, error) {
	return "PongPing", nil
}
