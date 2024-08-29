//go:generate veil
package main

import (
	"context"
	"fmt"
)

// @v:service
type UserService struct {
}

func (u *UserService) NewUser(ctx context.Context, name string) (string, error) {
	fmt.Println("NewUser")
	return "", nil
}

func (u *UserService) Say(ctx context.Context, userId string, msg string) (bool, error) {
	fmt.Println("Say")
	return true, nil
}
