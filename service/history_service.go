package service

import "adiubaidah/adi-bot/model/history"

type HistoryService interface {
	InsertHistory(senderId, recipientId, content, role string) error
	GetHistory(senderId string, recipientId string) (*[]history.History, error)
}
