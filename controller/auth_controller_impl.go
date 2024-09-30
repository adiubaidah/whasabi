package controller

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/middleware"
	"adiubaidah/adi-bot/model"
	"adiubaidah/adi-bot/service"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type AuthControllerImpl struct {
	AuthService service.AuthService
	service.UserService
}

func NewAuthController(authService service.AuthService, userService service.UserService) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
		UserService: userService,
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
	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Login success",
	})
}

func (controller *AuthControllerImpl) IsAuth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userContext := request.Context().Value(middleware.UserContext).(jwt.MapClaims)
	user := controller.UserService.Find(service.UserSearchParams{
		UserId: int(userContext["id"].(float64)),
	})

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   user,
	})
}

func (controller *AuthControllerImpl) Logout(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// Delete the token cookie
	cookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}

	http.SetCookie(writer, &cookie)

	helper.WriteToResponseBody(writer, &model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Logout success",
	})
}
