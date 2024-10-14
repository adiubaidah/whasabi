package controller

import (
	"net/http"
	"strconv"

	"github.com/adiubaidah/wasabi/app"
	"github.com/adiubaidah/wasabi/helper"
	"github.com/adiubaidah/wasabi/model"
	"github.com/adiubaidah/wasabi/service"

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
	user := controller.UserService.Create(request.Context(), userCreateRequest)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data: map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (controller *UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUpdateRequest := new(model.UserUpdateRequest)
	helper.ReadFromRequestBody(request, userUpdateRequest)

	userUpdateId := params.ByName("id")
	id, err := strconv.Atoi(userUpdateId)
	helper.PanicIfError("Error when convert", err)
	userUpdateRequest.ID = id

	user := controller.UserService.Update(request.Context(), userUpdateRequest)
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data: &model.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			IsActive: user.IsActive,
			Role:     user.Role,
		},
	})
}

func (controller *UserControllerImpl) List(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	users := controller.UserService.Find(&service.UserSearchParams{})
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   users,
	})
}

func (controller *UserControllerImpl) WebSocket(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(request.URL.Query().Get("id"))

	helper.PanicIfError("Error when convert", err)

	ai := controller.ProcessService.GetModel(uint(userId))

	controller.WsHub.ServeWebSocket(writer, request, ai.Phone)
}
