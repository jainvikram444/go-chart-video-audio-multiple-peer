package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Rooms for the strore all the room details with participates users
var Rooms Room

// CreateRoomRequest for the create room
func CreateRoomRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating Room Request")

	roomID := Rooms.CreateRoom()
	log.Println("Created room id : ", roomID)
	log.Println("Available room users: ", Rooms.RoomUsers)

	type respRoom struct {
		RoomID string `json:"room_id`
	}
	json.NewEncoder(w).Encode(respRoom{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type brodcastMessage struct {
	Message    map[string]interface{}
	RoomID     string
	ClientConn *websocket.Conn
}

var brodcast = make(chan brodcastMessage)

func brodcaster() {
	for {
		msg := <-brodcast

		for _, user := range Rooms.RoomUsers[msg.RoomID] {
			if user.Conn != msg.ClientConn {
				errCon := user.Conn.WriteJSON((msg.Message))

				if errCon != nil {
					log.Fatal(errCon)
					user.Conn.Close()
				}
			}
		}
	}
}

//JoinRoomRequest for the join room
func JoinRoomRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Join Room Request")

	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("Error: RoomID not foind in request")
		return
	}

	ws, errU := upgrader.Upgrade(w, r, nil)

	if errU != nil {
		log.Fatal("Error: Web Socket Upgrader: ", errU)
	}

	Rooms.InsertUserToRoom(roomID[0], false, ws)

	go brodcaster()

	for {
		var msg brodcastMessage

		err := ws.ReadJSON(msg.Message)

		if err != nil {
			log.Fatal("Error: Web socket read error: ", err)
		}

		msg.ClientConn = ws
		msg.RoomID = roomID[0]

		brodcast <- msg
	}
}
