//go:generate veil
package main

import (
	"context"
	"fmt"

	"github.com/hoyle1974/veil/veil"
)

type RoomId string

type Room struct {
	users map[UserId]any
}

func newRoom() *Room {
	return &Room{
		users: make(map[UserId]any),
	}
}

func (r *Room) addUser(userId UserId) {
	r.users[userId] = true
}

func (r *Room) removeUser(userId UserId) {
	delete(r.users, userId)
}

// @v:service -t rpc
type RoomService struct {
	rooms map[RoomId]*Room
}

func NewRoomService() *RoomService {
	return &RoomService{
		rooms: make(map[RoomId]*Room),
	}
}

func (r *RoomService) GetUsers(ctx context.Context, roomId RoomId) ([]UserId, error) {
	userids := []UserId{}

	for k, _ := range r.rooms[roomId].users {
		userids = append(userids, k)
	}
	return userids, nil
}

func (r *RoomService) AddUser(ctx context.Context, roomId RoomId, userId UserId) (int, bool, error) {
	fmt.Println("AddUser")

	room, ok := r.rooms[roomId]
	if !ok {
		room = newRoom()
		r.rooms[roomId] = room
	}

	room.addUser(userId)

	return 0, true, nil
}

func (r *RoomService) RemoveUser(ctx context.Context, roomId RoomId, userId UserId) (bool, error) {
	fmt.Println("RemoveUser")

	room, ok := r.rooms[roomId]
	if ok {
		room.removeUser(userId)
	}

	return true, nil
}

func (r *RoomService) Broadcast(ctx context.Context, roomId RoomId, msg string) (bool, error) {
	fmt.Println("Broadcast")

	users, err := veil.Lookup[UserService_Interface]()
	if err != nil {
		return false, err
	}

	room, ok := r.rooms[roomId]
	if ok {
		for k, _ := range room.users {
			users.Say(ctx, k, msg)
		}
	}

	return true, nil
}
