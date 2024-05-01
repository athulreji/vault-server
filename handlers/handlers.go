package handlers

import (
	"github.com/gorilla/websocket"
	"net/http"
	"signal_server/models"
	"signal_server/services"
	"signal_server/utils"
	"signal_server/utils/logger"
)

// For upgrading regular HTTP connections to WebSocket
var upgrader = websocket.Upgrader{}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade connection
	utils.CheckError(err)
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	// Read initial username
	var initMsg models.Message
	err = conn.ReadJSON(&initMsg)
	utils.CheckError(err)

	// Register client with username
	models.Clients[initMsg.From] = conn

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			logger.LogError(err.Error())
			// Remove client on error
			delete(models.Clients, initMsg.From)
			break
		}

		switch msg.Type {
		case "message":
			if msg.IsGroupMsg {
				services.ForwardGroupMsg(&msg)
			} else {
				services.ForwardDirectMsg(&msg)
			}
		case "cmd":
			switch msg.Content {
			case "create group":
				services.CreateGroup(&msg)
			case "join group":
				services.JoinGroup(&msg)
			}
		default:
			logger.LogWarning("Unknown Command")
		}

	}
}
