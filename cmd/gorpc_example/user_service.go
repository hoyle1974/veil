//go:generate veil
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// @v:service -t rpc
type UserService struct {
	users   map[string]string
	lastSay map[string]string
}

func NewUserService() *UserService {
	return &UserService{
		users:   make(map[string]string),
		lastSay: make(map[string]string),
	}
}

func (u *UserService) NewUser(ctx context.Context, name string) (string, error) {
	fmt.Println("NewUser")

	id, _ := uuid.NewRandom()
	u.users[id.String()] = name

	return id.String(), nil
}

func (u *UserService) Say(ctx context.Context, userId string, msg string) (bool, error) {
	u.lastSay[userId] = msg
	return true, nil
}

func (u *UserService) GetLastSay(ctx context.Context, userId string) (string, bool, error) {
	msg, ok := u.lastSay[userId]
	return msg, ok, nil
}
