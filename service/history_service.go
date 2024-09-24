package service

import "adiubaidah/adi-bot/model"

type HistoryService interface {
	InsertHistory(senderId, recipientId, content, role string) error
	GetHistory(senderId string, recipientId string) (*[]model.History, error)
}
