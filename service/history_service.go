package service

import "github.com/adiubaidah/wasabi/model"

type HistoryService interface {
	InsertHistory(processId uint, sender, receiver, content, role string) error
	GetHistory(sender string, receiver string) (*[]model.History, error)
}
