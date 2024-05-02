package main

import (
	"net/http"
	"signal_server/handlers"
	"signal_server/utils/logger"
)

func main() {
	logger.LogInfo("Server Listening on :8080")
	http.HandleFunc("/ws", handlers.HandleWebSocket)
	_ = http.ListenAndServe(":8080", nil)
}
