package app

import (
	"adiubaidah/adi-bot/helper"
	"context"
	"database/sql"
	"fmt"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func GetWaClient(phone string) *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	dsn := fmt.Sprintf("file:session/wa-%s.db?_foreign_keys=on", phone)

	container, err := sqlstore.New("sqlite3", dsn, dbLog)
	helper.PanicIfError("", err)

	db, err := sql.Open("sqlite3", "file:history.db?_foreign_keys=on")
	helper.PanicIfError("Error opening database:", err)
	defer db.Close()

	device, err := container.GetFirstDevice()
	helper.PanicIfError("", err)
	client := whatsmeow.NewClient(device, dbLog)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		helper.PanicIfError("", err)

		for evt := range qrChan {
			if evt.Event == "code" {
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "public/qr-"+phone+".png")
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
		helper.PanicIfError("", err)
	}

	return client
}
