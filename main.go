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
	Group      string `json:"group"`
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
		} else if msg.Type == "cmd" {
			if msg.Content == "create group" {
				createGroup(&msg)
			} else if msg.Content == "join group" {
				joinGroup(&msg)
			}
		}
	}
}

func createGroup(msg *Message) {
	if _, ok := groups[msg.Group]; ok {
		println("Group already present.")
	} else {
		groups[msg.Group] = append(groups[msg.Group], msg.From)
	}
}

func joinGroup(msg *Message) {
	if _, ok := groups[msg.Group]; !ok {
		println("Group doesn't exist: ", msg.Group)
	} else {
		groups[msg.Group] = append(groups[msg.Group], msg.From)
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
	if _, ok := groups[msg.Group]; !ok {
		fmt.Println("Group doesn't exist: ", msg.Group)
		return
	}
	group := groups[msg.Group]
	flag := false
	for _, i := range group {
		if i == msg.From {
			flag = true
		}
	}
	if !flag {
		fmt.Println(msg.From, " not belong to Group: ", msg.Group)
		return
	}

	for _, client := range groups[msg.Group] {
		if client != msg.From { // Skip sending to self
			err := clients[client].WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.ListenAndServe(":8080", nil)
}
