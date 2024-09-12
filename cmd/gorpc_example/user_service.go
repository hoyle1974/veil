//go:generate veil
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserId string

// @v:service -t rpc
type UserService struct {
	users   map[UserId]string
	lastSay map[UserId]string
}

func NewUserService() *UserService {
	return &UserService{
		users:   make(map[UserId]string),
		lastSay: make(map[UserId]string),
	}
}

func (u *UserService) NewUser(ctx context.Context, name string) (UserId, error) {
	fmt.Println("NewUser")

	id, _ := uuid.NewRandom()
	u.users[UserId(id.String())] = name

	return UserId(id.String()), nil
}

func (u *UserService) Say(ctx context.Context, userId UserId, msg string) (bool, error) {
	u.lastSay[userId] = msg
	return true, nil
}

func (u *UserService) GetLastSay(ctx context.Context, userId UserId) (string, bool, error) {
	msg, ok := u.lastSay[userId]
	return msg, ok, nil
}
