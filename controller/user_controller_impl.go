package controller

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService    service.UserService
	ProcessService service.ProcessService
	WsHub          *app.WebSocketHub
	Validate       *validator.Validate
}

func NewUserController(userService service.UserService, processService service.ProcessService, wsHub *app.WebSocketHub, validate *validator.Validate) UserController {
	return &UserControllerImpl{
		UserService:    userService,
		ProcessService: processService,
		WsHub:          wsHub,
		Validate:       validate,
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

	helper.PanicIfError("Error when convert", err)

	ai := controller.ProcessService.GetModel(uint(userId))

	controller.WsHub.ServeWebSocket(writer, request, ai.Phone)
}
