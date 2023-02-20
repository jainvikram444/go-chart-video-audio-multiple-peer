package main

import (
	"fmt"
	"go-chart-video-audio-multiple-peer/server"
	"log"
	"net/http"
)

func main() {
	const (
		httpPort string = ":3000"
	)
	server.Rooms.Init()

	fmt.Println("Welcome to WebRTC peer to peeer video call between the muliple users.")

	http.HandleFunc("/create_room", server.CreateRoomRequest)
	http.HandleFunc("/join_room", server.JoinRoomRequest)

	log.Println("Initialize the Server")

	log.Printf("Server successfully running on port: %v", httpPort)
	errHttp := http.ListenAndServe(httpPort, nil)
	if errHttp != nil {
		fmt.Println("Error: some issue in starting server: ", errHttp)
		log.Fatal(errHttp)
		return
	}
	log.Printf("Server successfully stop on port: %v", httpPort)
}
