package service

import (
	"adiubaidah/adi-bot/model/history"

	"gorm.io/gorm"
)

type HistoryServiceImpl struct {
	db *gorm.DB
}

func NewHistoryService(db *gorm.DB) HistoryService {
	return &HistoryServiceImpl{db: db}
}

func (s *HistoryServiceImpl) InsertHistory(senderId, receiver, content, role string) error {
	history := history.History{
		Sender:   senderId,
		Receiver: receiver,
		Content:  content,
		RoleAs:   role,
	}
	return s.db.Create(&history).Error
}

func (s *HistoryServiceImpl) GetHistory(senderId string, recipientId string) (*[]history.History, error) {
	var histories []history.History
	err := s.db.Where("(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)", senderId, recipientId, recipientId, senderId).
		Order("created_at ASC").
		Find(&histories).Error
	return &histories, err
}
