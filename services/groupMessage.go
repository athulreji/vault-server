package services

import (
	"signal_server/models"
	"signal_server/utils/logger"
)

func CreateGroup(msg *models.Message) {
	if _, ok := models.Groups[msg.Group]; ok {
		logger.LogWarning("Group already exist: " + msg.Group)
		return
	}
	models.Groups[msg.Group] = append(models.Groups[msg.Group], msg.From)
}

func JoinGroup(msg *models.Message) {
	if _, ok := models.Groups[msg.Group]; !ok {
		logger.LogWarning("Group does not exist: " + msg.Group)
		return
	}
	models.Groups[msg.Group] = append(models.Groups[msg.Group], msg.From)
}

func ForwardGroupMsg(msg *models.Message) {
	if _, ok := models.Groups[msg.Group]; !ok {
		logger.LogWarning("Group does not exist: " + msg.Group)
		return
	}
	group := models.Groups[msg.Group]
	flag := false
	for _, i := range group {
		if i == msg.From {
			flag = true
		}
	}

	if !flag {
		logger.LogWarning(msg.From + " not belong to Group: " + msg.Group)
		return
	}

	for _, client := range models.Groups[msg.Group] {
		if client != msg.From {
			// Skip sending to self
			err := models.Clients[client].WriteJSON(msg)
			if err != nil {
				logger.LogError(err.Error())
				delete(models.Clients, client)
			}
		}
	}
}
