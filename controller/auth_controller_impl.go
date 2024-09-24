package controller

import (
	"adiubaidah/adi-bot/exception"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type AuthControllerImpl struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
	}
}

func (controller *AuthControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userLoginRequest := new(model.UserLoginRequest)
	helper.ReadFromRequestBody(request, userLoginRequest)

	token := controller.AuthService.Login(request.Context(), *userLoginRequest)

	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 72),
	}
	http.SetCookie(writer, &cookie)
	helper.WriteToResponseBody(writer, "Bearer "+token)

}

func (controller *AuthControllerImpl) Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userLogoutRequest := new(model.UserLogoutRequest)
	helper.ReadFromRequestBody(request, userLogoutRequest)
	err := controller.AuthService.Logout(request.Context(), *userLogoutRequest)
	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}
	helper.WriteToResponseBody(writer, "Logout successful")
}
