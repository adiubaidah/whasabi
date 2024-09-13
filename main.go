package main

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/controller"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/routes"
	"adiubaidah/adi-bot/service"
	"context"
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app.Init()
	validate := validator.New()
	db, err := sql.Open("sqlite3", "file:history.db?_foreign_keys=on")
	helper.PanicIfError("Error opening database:", err)
	defer db.Close()

	context := context.Background()

	aiClient := app.GetAIClient(context)
	waService := service.NewWaService()
	aiService := service.NewAiService(aiClient, validate)
	// authService := service.NewAuthService(db, validate)
	historyService := service.NewHistoryService(db)

	aiController := controller.NewAiController(waService, aiService, historyService, validate)

	router := routes.SetupRouter(aiController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
	err = server.ListenAndServe()
	helper.PanicIfError("failed to start server", err)

}
