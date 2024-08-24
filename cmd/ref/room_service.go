//go:generate veil
package main

import (
	"context"
	"fmt"
)

// @d:service
type RoomService struct {
}

func (r *RoomService) AddUser(ctx context.Context, roomId string, userId string) (bool, error) {
	fmt.Println("AddUser")
	return true, nil
}

func (r *RoomService) RemoveUser(ctx context.Context, roomId string, userId string) (bool, error) {
	fmt.Println("RemoveUser")
	return true, nil
}

func (r *RoomService) Broadcast(ctx context.Context, roomId string, msg string) (bool, error) {
	fmt.Println("Broadcast")

	return true, nil
}
