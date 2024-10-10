package app

import (
	"fmt"

	"github.com/adiubaidah/wasabi/helper"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func GetWaClient(phone string) (*whatsmeow.Client, *sqlstore.Container) {
	dbLog := waLog.Stdout("Database", "ERROR", true)

	dsn := fmt.Sprintf("file:session/wa-%s.db?_foreign_keys=on", phone)

	container, err := sqlstore.New("sqlite3", dsn, dbLog)
	helper.PanicIfError("", err)

	device, err := container.GetFirstDevice()
	helper.PanicIfError("", err)
	client := whatsmeow.NewClient(device, dbLog)
	return client, container
}
