package services

import (
	"signal_server/models"
	"signal_server/utils/logger"
)

func ForwardDirectMsg(msg *models.Message) {
	if _, ok := models.Clients[msg.To]; ok {
		err := models.Clients[msg.To].WriteJSON(msg)
		if err != nil {
			logger.LogError(err.Error())
			delete(models.Clients, msg.To)
		}
	}
}
