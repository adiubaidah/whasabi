package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/generative-ai-go/genai"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	client    *whatsmeow.Client
	startTime time.Time
)

func handleIncomingMessage(message *events.Message, db *sql.DB) {

	context := context.Background()
	// Ambil teks pesan
	text := message.Message.GetConversation()
	fmt.Println("Pesan masuk:", text)
	senderId := message.Info.Sender.String()

	aiResponse, clientJID, err := getAiReponse(db, senderId, text)
	PanicIfError("Error when getting AI response:", err)
	err = InsertHistory(db, senderId, clientJID, text, "user")
	PanicIfError("Error saat menyimpan riwayat:", err)
	err = InsertHistory(db, clientJID, senderId, aiResponse, "model")
	PanicIfError("Error saat menyimpan riwayat:", err)

	_, err = client.SendMessage(context, message.Info.Chat, &waE2E.Message{
		Conversation: proto.String(aiResponse),
	})
	if err != nil {
		fmt.Println("Error saat mengirim pesan:", err)
	}
}

func getAiReponse(db *sql.DB, sender string, input string) (string, string, error) {
	ctx := context.Background()
	model, aiClient := GetAIModel()
	defer aiClient.Close() // Tutup klien AI setelah selesai mengirim pesan agar tidak membuang sumber daya

	// session.
	session := model.StartChat()
	var sessionHistory []*genai.Content

	getHistory, err := GetHistory(db, sender, client.Store.ID.String())
	PanicIfError("Error saat mengambil riwayat:", err)
	for _, history := range *getHistory {
		sessionHistory = append(sessionHistory, &genai.Content{
			Role:  history.Role,
			Parts: []genai.Part{genai.Text(history.Content)},
		})
	}
	if len(sessionHistory) > 0 {
		session.History = sessionHistory
	}

	resp, err := session.SendMessage(ctx, genai.Text(input))
	PanicIfError("Error saat mengambil respon:", err)
	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), client.Store.ID.String(), nil

}

func main() {
	if err := InitConfig(); err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:log.db?_foreign_keys=on", dbLog)
	PanicIfError("", err)

	// Buat koneksi ke SQLite
	db, err := sql.Open("sqlite3", "file:history.db?_foreign_keys=on")
	PanicIfError("Error opening database:", err)
	defer db.Close()

	// Buat tabel history jika belum ada
	err = CreateHistoryTableIfNotExists(db)
	PanicIfError("Error saat membuat tabel history:", err)

	StartHistoryCleanup(db, 1*time.Hour)

	// Inisialisasi WhatsMeow client
	device, err := container.GetFirstDevice()
	PanicIfError("", err)

	client = whatsmeow.NewClient(device, dbLog)
	client.AddEventHandler(func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			if v.Info.IsGroup {
				return
			}

			if v.Info.IsFromMe {
				return
			}

			fmt.Println("Pesan timestamp:", v.Info.Timestamp)

			if v.Info.Timestamp.After(startTime) {
				handleIncomingMessage(v, db)
			}
		default:
			fmt.Println("Event lainnya:", evt)
		}
	})

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		PanicIfError("", err)

		for evt := range qrChan {
			if evt.Event == "code" {
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qr.png")
				if err != nil {
					fmt.Println("Failed to generate QR code:", err)
				} else {
					fmt.Println("QR code generated and saved as qr.png. Scan it using WhatsApp!")
				}
			} else {
				fmt.Println("QR event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		PanicIfError("", err)
	}
	// Catat waktu saat bot terhubung
	startTime = time.Now()

	fmt.Println("Start time:", startTime)

	// Menangani sinyal berhenti
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Menjaga aplikasi tetap berjalan
	fmt.Println("Aplikasi berjalan. Tekan Ctrl+C untuk berhenti.")
	select {
	case <-stop:
		// db.Close()
		fmt.Println("Aplikasi berhenti.")
	}

}
