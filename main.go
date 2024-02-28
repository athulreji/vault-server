package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // For upgrading regular HTTP connections to WebSocket

type Message struct {
	Type       string `json:"type"`
	IsGroupMsg bool   `json:"isGroupMsg"`
	To         string `json:"to"`
	Content    string `json:"content"`
	From       string `json:"From"`
}

var clients = make(map[string]*websocket.Conn) // To store client connections with usernames

var groups = make(map[string][]string)

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

	clients[initMsg.From] = conn // Register client with username

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, initMsg.From) // Remove client on error
			break
		}

		if msg.Type == "message" {
			if msg.IsGroupMsg {
				forwardGroupMsg(&msg)
			} else {
				forwardDirectMsg(&msg)
			}
		}
	}
}

func forwardDirectMsg(msg *Message) {
	if _, ok := clients[msg.To]; ok {
		err := clients[msg.To].WriteJSON(msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, msg.To)
		}
	}
}

func forwardGroupMsg(msg *Message) {
	// Broadcast to all other clients
	// for client := range clients {
	// 	if client != conn { // Skip sending to self
	// 		newMsg := Message{
	// 			Type:    "message",
	// 			From:    msg.From,
	// 			Content: msg.Content,
	// 		}
	// 		err := client.WriteJSON(newMsg)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			delete(clients, client)
	// 		}
	// 	}
	// }
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.ListenAndServe(":8080", nil)
}
