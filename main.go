package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // For upgrading regular HTTP connections to WebSocket

type Message struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Content string `json:"content"`
	From    string `json:"From"`
}

var clients = make(map[*websocket.Conn]string) // To store client connections with usernames

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade connection
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Read initial username
	var initMsg Message
	err = conn.ReadJSON(&initMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	clients[conn] = initMsg.From // Register client with username

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn) // Remove client on error
			break
		}

		// Broadcast to all other clients
		for client := range clients {
			if client != conn { // Skip sending to self
				newMsg := Message{
					Type:    "message",
					From:    msg.From,
					Content: msg.Content,
				}
				err := client.WriteJSON(newMsg)
				if err != nil {
					fmt.Println(err)
					delete(clients, client)
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.ListenAndServe(":8080", nil)
}
