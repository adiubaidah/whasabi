package main

import (
	"database/sql"
	"fmt"
	"time"
)

type History struct {
	ID        int
	Sender    string
	Recipient string
	Content   string
	Role      string
	CreatedAt time.Time
}

func CreateHistoryTableIfNotExists(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS history (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender TEXT NOT NULL,
        recipient TEXT NOT NULL,
        content TEXT NOT NULL,
        role TEXT CHECK( role IN ('user', 'model') ) NOT NULL,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(query)
	return err
}

func InsertHistory(db *sql.DB, senderId, recipientId, content, role string) error {
	query := `INSERT INTO history (sender, recipient,  content, role) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, senderId, recipientId, content, role)
	return err
}

func GetHistory(db *sql.DB, senderId string, recipientId string) (*[]History, error) {
	query := `SELECT id, sender, recipient, content, role, createdAt
    FROM history  
    WHERE sender = ? 
       OR recipient = ? OR sender = ? OR recipient = ?
    ORDER BY createdAt ASC;`
	rows, err := db.Query(query, senderId, recipientId, senderId, recipientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []History
	for rows.Next() {
		var history History
		err = rows.Scan(&history.ID, &history.Sender, &history.Recipient, &history.Content, &history.Role, &history.CreatedAt)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return &histories, nil
}

func DeleteOldHistory(db *sql.DB) error {
	query := `DELETE FROM history WHERE createdAt < datetime('now', '-1 hour')`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("Old history records deleted")
	return nil
}

func StartHistoryCleanup(db *sql.DB, interval time.Duration) {
	go func() {
		for {
			err := DeleteOldHistory(db)
			PanicIfError("Error saat menghapus riwayat lama:", err)
			time.Sleep(interval)
		}
	}()
}
