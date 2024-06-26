package chat

import (
	"fmt"
	"net/http"
	"streamvault/postgres"
	"strings"

	"github.com/gorilla/websocket"
)

// type Connection struct {
// 	conn *websocket.Conn
// }

var clients = map[string]map[*websocket.Conn]bool{}
var Broadcast = make(chan Message)

	type Message struct {
		User   postgres.UserDetails `json:"user"`
		Message  string  `json:"message"`
		StreamId string  `json:"streamId"`
	}

type ErrorResponse struct {
	Error string `json:"error"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"streamId"}, // <-- add this line
}

func Chat(w http.ResponseWriter, r *http.Request) {
	var a string = r.Header.Get("Sec-WebSocket-Protocol")
	parts := strings.Split(a, ", ")
	if len(parts) < 3 {
		fmt.Println("Invalid Sec-WebSocket-Protocol header")
		return
	}
	var streamId string = parts[1]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
		return
	}

	fmt.Println("/" + streamId + "/")
	if _, exists := clients[streamId]; !exists {
		clients[streamId] = map[*websocket.Conn]bool{}
	}
	clients[streamId][conn] = true

	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("errror",err)
			clients[streamId][conn] = false
			delete(clients[streamId], conn)
			fmt.Println(clients,"clients")
			return
		}
		if msg.User.UserId == "" {
			conn.WriteJSON(ErrorResponse{Error: "Log in to send messages"})
			continue
		}
		exists, err := postgres.UserExists(msg.User.UserId)

		if err != nil {
			fmt.Println(err)
			conn.WriteJSON(ErrorResponse{Error: "Error checking if user exists"})
			continue
		}
		if !exists {

			conn.WriteJSON(ErrorResponse{Error: "User does not exist"})
			continue

		}

		fmt.Println(msg.Message, msg.User.UserId, msg.StreamId)

		
		

		



		Broadcast <- msg
	}

}

func HandleMessages() {
	for {
		msg := <-Broadcast
		fmt.Println("Received a message: ",msg.Message)
		var streamId string = msg.StreamId
		err:=postgres.PostChat(msg.StreamId, msg.User.UserId, msg.Message)
		if err != nil {
			fmt.Println(err)
			
			continue
		}
		for client := range clients[streamId] {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients[msg.StreamId], client)
			}
		}
	}
}
