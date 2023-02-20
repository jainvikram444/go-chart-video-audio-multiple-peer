package server

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Host bool
	Conn *websocket.Conn
}

type Room struct {
	Mutex     sync.RWMutex
	RoomUsers map[string][]User
}

// Init for the initialize the room
func (r *Room) Init() {
	r.RoomUsers = make(map[string][]User)
}

// Get for the fetch user details of the request room
func (r *Room) Get(roomID string) []User {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	return r.RoomUsers[roomID]
}

// CreateRoom for the create a new room
func (r *Room) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	var stringKey = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	byteKey := make([]rune, 8)

	for i := range byteKey {
		byteKey[i] = stringKey[rand.Intn(len(stringKey))]
	}

	roomID := string(byteKey)
	r.RoomUsers[roomID] = []User{}

	return roomID
}

// InsertUserToRoom  for the insert new user to room
func (r *Room) InsertUserToRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	user := User{host, conn}
	r.RoomUsers[roomID] = append(r.RoomUsers[roomID], user)

	fmt.Printf("Insert user to room: %v", roomID)
}

// DeleteRoom for the delete room
func (r *Room) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.RoomUsers, roomID)
}
