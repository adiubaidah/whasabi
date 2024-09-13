package controller

import (
	"adiubaidah/adi-bot/service"
	"net/http"
)

type AuthControllerImpl struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
	}
}

func (controller *AuthControllerImpl) Login(writer http.ResponseWriter, request *http.Request) {

}

func (controller *AuthControllerImpl) Logout(writer http.ResponseWriter, request *http.Request) {

}
