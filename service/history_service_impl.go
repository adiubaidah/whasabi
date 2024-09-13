package service

import (
	"adiubaidah/adi-bot/model/history"
	"database/sql"
	"fmt"
	"time"
)

type HistoryServiceImpl struct {
	db *sql.DB
}

func NewHistoryService(db *sql.DB) HistoryService {
	return &HistoryServiceImpl{db: db}
}

func (s *HistoryServiceImpl) InsertHistory(senderId, recipientId, content, role string) error {
	query := `INSERT INTO history (sender, recipient, content, role) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, senderId, recipientId, content, role)
	return err
}

func (s *HistoryServiceImpl) GetHistory(senderId string, recipientId string) (*[]history.History, error) {
	query := `SELECT id, sender, recipient, content, role, createdAt
    FROM history  
    WHERE (sender = ? AND recipient = ?) 
       OR (sender = ? AND recipient = ?)
    ORDER BY createdAt ASC;`
	rows, err := s.db.Query(query, senderId, recipientId, recipientId, senderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []history.History
	for rows.Next() {
		var history history.History
		err = rows.Scan(&history.ID, &history.Sender, &history.Recipient, &history.Content, &history.Role, &history.CreatedAt)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return &histories, nil
}

func (s *HistoryServiceImpl) DeleteOldHistory() error {
	query := `DELETE FROM history WHERE createdAt < datetime('now', '-1 hour')`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("Old history records deleted")
	return nil
}

func (s *HistoryServiceImpl) StartHistoryCleanup(interval time.Duration) {
	go func() {
		for {
			err := s.DeleteOldHistory()
			if err != nil {
				fmt.Println("Error deleting old history records:", err)
			}
			time.Sleep(interval)
		}
	}()
}
