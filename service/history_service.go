package service

import "adiubaidah/adi-bot/model"

type HistoryService interface {
	InsertHistory(sender, receiver, content, role string) error
	GetHistory(sender string, receiver string) (*[]model.History, error)
}
