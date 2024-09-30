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
	waService := service.NewWaService(websocketHub, db)
	aiService := service.NewAiService(aiClient, db, validate)
	authService := service.NewAuthService(db, validate)
	historyService := service.NewHistoryService(db)
	userService := service.NewUserService(db, validate)

	aiController := controller.NewAiController(waService, aiService, historyService, validate)
	authController := controller.NewAuthController(authService, userService)
	userController := controller.NewUserController(userService, websocketHub)

	router := routes.SetupRouter(aiController, userController, authController)

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := helper.GetEnv("APP_ORIGIN")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				fmt.Println("OPTIONS")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	PORT := helper.GetEnv("APP_PORT")
	server := http.Server{
		Addr:    "localhost:" + PORT,
		Handler: corsHandler(router),
	}
	fmt.Println("Server started at localhost:" + PORT)
	err := server.ListenAndServe()
	helper.PanicIfError("failed to start server", err)

}
