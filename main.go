package main

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/controller"
	"adiubaidah/adi-bot/db"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/routes"
	"adiubaidah/adi-bot/service"
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func main() {
	app.Init()
	validate := validator.New()
	db := db.NewDB()

	context := context.Background()

	aiClient := app.GetAIClient(context)
	websocketHub := app.WsHub
	waService := service.NewWaService(websocketHub)
	aiService := service.NewAiService(aiClient, db, validate)
	authService := service.NewAuthService(db, validate)
	historyService := service.NewHistoryService(db)
	userService := service.NewUserService(db, validate)

	aiController := controller.NewAiController(waService, aiService, historyService, validate)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService, websocketHub)

	router := routes.SetupRouter(aiController, userController, authController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
	fmt.Println("Server started at localhost:3000")
	err := server.ListenAndServe()
	helper.PanicIfError("failed to start server", err)

}
