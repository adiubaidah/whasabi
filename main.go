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
	"os"

	"github.com/go-playground/validator/v10"
)

func main() {
	validate := validator.New()
	db := db.NewDB()

	context := context.Background()

	aiClient := app.GetAIClient(context)
	websocketHub := app.WsHub
	processService := service.NewProcessService(aiClient, websocketHub, db)
	authService := service.NewAuthService(db)
	historyService := service.NewHistoryService(db)
	userService := service.NewUserService(db)

	processController := controller.NewProcessController(processService, historyService, validate)
	authController := controller.NewAuthController(authService, userService, validate)
	userController := controller.NewUserController(userService, processService, websocketHub, validate)

	router := routes.SetupRouter(processController, userController, authController)

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := os.Getenv("APP_ORIGIN")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	PORT := os.Getenv("APP_PORT")
	server := http.Server{
		Addr:    "localhost:" + PORT,
		Handler: corsHandler(router),
	}
	fmt.Println("Server started at localhost:" + PORT)
	err := server.ListenAndServe()
	helper.PanicIfError("failed to start server", err)

}
