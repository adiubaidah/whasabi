package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/adiubaidah/wasabi/app"
	"github.com/adiubaidah/wasabi/controller"
	"github.com/adiubaidah/wasabi/db"
	"github.com/adiubaidah/wasabi/helper"
	"github.com/adiubaidah/wasabi/model"
	"github.com/adiubaidah/wasabi/routes"
	"github.com/adiubaidah/wasabi/service"

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
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-ticker.C
			model.DeleteOldHistories(db, 3)
		}
	}()
	PORT := os.Getenv("APP_PORT")
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), PORT),
		Handler: corsHandler(router),
	}
	fmt.Println("Server has started at localhost:" + PORT)
	err := server.ListenAndServe()
	helper.PanicIfError("failed to start server", err)

}
