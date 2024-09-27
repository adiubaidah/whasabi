package controller

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
	WsHub       *app.WebSocketHub
}

func NewUserController(userService service.UserService, wsHub *app.WebSocketHub) UserController {
	return &UserControllerImpl{
		UserService: userService,
		WsHub:       wsHub,
	}
}

func (controller *UserControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userCreateRequest := new(model.UserCreateRequest)
	helper.ReadFromRequestBody(request, userCreateRequest)
	user := controller.UserService.Create(request.Context(), *userCreateRequest)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data: map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (controller *UserControllerImpl) WebSocket(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(request.URL.Query().Get("id"))

	helper.PanicIfError("", err)

	ai := controller.UserService.GetService(userId)

	controller.WsHub.ServeWebSocket(writer, request, ai.Phone)
}
