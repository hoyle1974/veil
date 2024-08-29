//go:generate veil
package main

import (
	"context"
	"fmt"
	"time"
)

// @d:service
type RoomService struct {
}

func (r *RoomService) AddUser(ctx context.Context, roomId string, userId string) (int, bool, error) {
	fmt.Println("AddUser")
	return 0, true, nil
}

func (r *RoomService) RemoveUser(ctx context.Context, roomId string, userId string) (bool, error) {
	fmt.Println("RemoveUser")
	return true, nil
}

func (r *RoomService) Broadcast(ctx context.Context, roomId string, msg string) (bool, error) {
	fmt.Println("Broadcast - start")

	time.Sleep(time.Duration(5) * time.Second)

	fmt.Println("Broadcast - end")

	return true, nil
}
