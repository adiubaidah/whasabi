package service

import (
	"adiubaidah/adi-bot/model"

	"gorm.io/gorm"
)

type HistoryServiceImpl struct {
	DB *gorm.DB
}

func NewHistoryService(DB *gorm.DB) HistoryService {
	return &HistoryServiceImpl{DB: DB}
}

func (s *HistoryServiceImpl) InsertHistory(senderId, receiver, content, role string) error {
	history := model.History{
		Sender:   senderId,
		Receiver: receiver,
		Content:  content,
		RoleAs:   role,
	}
	return s.DB.Create(&history).Error
}

func (s *HistoryServiceImpl) GetHistory(senderId string, recipientId string) (*[]model.History, error) {
	var histories []model.History
	err := s.DB.Where("(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)", senderId, recipientId, recipientId, senderId).
		Order("created_at ASC").
		Find(&histories).Error
	return &histories, err
}
