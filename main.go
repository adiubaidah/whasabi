package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/generative-ai-go/genai"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

var (
	client     *whatsmeow.Client
	statusChan = make(chan bool)
)

func eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsGroup {
			return
		}
		handleIncomingMessage(v)
	case *events.Presence:
		if v.Unavailable {
			statusChan <- false
		} else {
			statusChan <- true
		}

	default:
		fmt.Println("Event lainnya:", evt)
	}
}

func handleIncomingMessage(message *events.Message) {

	context := context.Background()
	// Ambil teks pesan
	text := message.Message.GetConversation()
	aiResponse, err := getAiReponse(text)
	if err != nil {
		fmt.Println("Error saat mengirim pesan:", err)
		return
	}

	_, err = client.SendMessage(context, message.Info.Chat, &waE2E.Message{
		Conversation: proto.String(aiResponse),
	})
	if err != nil {
		fmt.Println("Error saat mengirim pesan:", err)
	}
}

func getAiReponse(input string) (string, error) {
	ctx := context.Background()
	apiKey, err := GetConfig("GENAI_API_KEY")
	if err != nil {
		fmt.Println("Error saat mengambil api key:", err)
		return "", err
	}

	modelID, err := GetConfig("GENAI_MODEL_ID")
	if err != nil {
		fmt.Println("Error saat mengambil model:", err)
		return "", err
	}

	systemInstruction, err := GetConfig("GENAI_SYSTEM_INSTRUCTION")
	if err != nil {
		fmt.Println("Error saat memberi instruksi:", err)
		return "", err
	}

	aiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer aiClient.Close()

	model := aiClient.GenerativeModel(modelID)
	model.SystemInstruction = genai.NewUserContent(genai.Text(systemInstruction))

	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		fmt.Println("Error saat mengambil respon:", err)
		return "", err
	}
	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), nil

}

func main() {
	if err := InitConfig(); err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:log.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	// Inisialisasi WhatsMeow client
	device, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	client = whatsmeow.NewClient(device, dbLog)

	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}

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
		if err != nil {
			panic(err)
		}
	}
	// Menangani sinyal berhenti
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Menjaga aplikasi tetap berjalan
	fmt.Println("Aplikasi berjalan. Tekan Ctrl+C untuk berhenti.")
	select {
	case <-stop:
		fmt.Println("Aplikasi berhenti.")
	}

}
