package models

import "github.com/gorilla/websocket"

// To store client connections with usernames
var Clients = make(map[string]*websocket.Conn)
var Groups = make(map[string][]string)

type Message struct {
	Type       string `json:"type"`
	IsGroupMsg bool   `json:"isGroupMsg"`
	Group      string `json:"group"`
	To         string `json:"to"`
	Content    string `json:"content"`
	From       string `json:"from"`
}
