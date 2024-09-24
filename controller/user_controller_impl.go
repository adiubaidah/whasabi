package controller

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userCreateRequest := new(model.UserCreateRequest)
	helper.ReadFromRequestBody(request, userCreateRequest)
	user := controller.UserService.Create(request.Context(), *userCreateRequest)
	helper.WriteToResponseBody(writer, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}
